package sevsnp

// The same test for every variant of the scheme: what each one answers to the
// same real reports, and how many times it asks AMD's KDS. Not for submission.

import (
	_ "embed"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-sev-guest/abi"
	"github.com/google/go-sev-guest/verify/trust"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/veraison/corim/comid"
	"github.com/veraison/corim/corim"
	"github.com/veraison/ratsd/tokens"
	"github.com/veraison/services/vts/appraisal"
)

//go:embed test/evidence/sevsnp-aws-vlek-tsm-report.cbor
var mxAWS []byte

//go:embed test/evidence/sevsnp-gcp-tsm-report.cbor
var mxGCP []byte

const (
	mxAWSReportData = "6131f96c32e9cb962b4c0c7507c9c7fe36ec4c527a73a7809c786251dbe5809fa9b84cd8aec306f3d151a78acb323d8e753e085aa6f829d636792a1ae43fba44"
	mxGCPReportData = "f896ba422824bb5bb0201da4599c9124d24ab896cc3cec0a0c3a9ceac9a5208f4cbc7d3aa0f895ccc61cf7811c5b2cc62ea41c70771150e64c09cb7668539524"
	// the nonce of upstream's own sample (scheme_test.go, Test_ValidateEvidenceIntegrity_ok)
	mxUpstreamNonce  = "4d4944424e48323869696f69736a5079787878787878787878787878787878784d4944424e48323869696f69736a507978787878787878787878787878787878"
	mxRatsdMediaType = `application/eat+cwt; eat_profile="tag:github.com,2025:veraison/ratsd/cmw"`
)

type mxTransport struct {
	base http.RoundTripper
	urls []string
}

func (c *mxTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if strings.Contains(r.URL.Host, "kdsintf.amd.com") {
		c.urls = append(c.urls, r.URL.Path)
	}
	return c.base.RoundTrip(r)
}

// mxLeafOnly keeps the one entry of the table that holds the report's signing
// certificate (VCEK or VLEK): the table of a host that supplies no chain.
func mxLeafOnly(t *testing.T, table []byte) []byte {
	t.Helper()
	entries, err := abi.ParseSnpCertTableHeader(table)
	require.NoError(t, err)
	for _, e := range entries {
		g := e.GUID.String()
		if g == abi.VcekGUID || g == abi.VlekGUID {
			der := table[e.Offset : e.Offset+e.Length]
			ct := abi.CertTable{Entries: []abi.CertTableEntry{{GUID: uuid.MustParse(g), RawCert: der}}}
			return ct.Marshal()
		}
	}
	t.Fatal("no VCEK/VLEK in table")
	return nil
}

func mxToken(t *testing.T, src []byte, table []byte) []byte {
	t.Helper()
	var tsm tokens.TSMReport
	require.NoError(t, tsm.FromCBOR(src))
	if table != nil {
		tsm.AuxBlob = table
	}
	out, err := tsm.ToCBOR()
	require.NoError(t, err)
	return out
}

func mxAMDAnchor(t *testing.T, line string) []*comid.KeyTriple {
	t.Helper()
	roots := trust.DefaultRootCerts[line]
	require.NotNil(t, roots)
	arkPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: roots.ProductCerts.Ark.Raw})
	ta, err := comid.NewPKIXBase64Cert(string(arkPEM))
	require.NoError(t, err)
	return []*comid.KeyTriple{{VerifKeys: comid.CryptoKeys{ta}}}
}

func mxUpstreamAnchor(t *testing.T) []*comid.KeyTriple {
	t.Helper()
	var endCorim corim.UnsignedCorim
	var endComid comid.Comid
	require.NoError(t, endCorim.FromCBOR(sevsnpCorimValidEndorsements))
	require.NoError(t, endComid.FromCBOR(endCorim.Tags[1].Content))
	keyTriples := *endComid.Triples.AttestVerifKeys
	return []*comid.KeyTriple{&keyTriples[0]}
}

var mxDetail = regexp.MustCompile(`"detail":\["([^"\\]*)`)

func mxShort(err error) string {
	if err == nil {
		return "ok"
	}
	s := strings.ReplaceAll(err.Error(), `\"`, `"`)
	if all := mxDetail.FindAllStringSubmatch(s, -1); all != nil {
		s = all[len(all)-1][1] // the innermost detail of nested BadEvidence errors
	}
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > 90 {
		s = s[:90] + "…"
	}
	return "FAIL: " + s
}

func Test_Matrix(t *testing.T) {
	variant, gsg := os.Getenv("MX_VARIANT"), os.Getenv("MX_GSG")
	ct := &mxTransport{base: http.DefaultTransport}
	http.DefaultTransport = ct
	defer func() { http.DefaultTransport = ct.base }()

	awsData, err := hex.DecodeString(mxAWSReportData)
	require.NoError(t, err)
	gcpData, err := hex.DecodeString(mxGCPReportData)
	require.NoError(t, err)
	upNonce, err := hex.DecodeString(mxUpstreamNonce)
	require.NoError(t, err)

	var awsTSM, gcpTSM tokens.TSMReport
	require.NoError(t, awsTSM.FromCBOR(mxAWS))
	require.NoError(t, gcpTSM.FromCBOR(mxGCP))

	type row struct {
		name      string
		data      []byte
		mediaType string
		nonce     []byte
		ta        []*comid.KeyTriple
		keepCache bool // run right after the previous row, without clearing go-sev-guest's product-chain cache
	}
	rows := []row{
		{"upstream's own sample: Genoa VCEK, PEM table, ratsd token", sevsnpRatsdToken, mxRatsdMediaType, upNonce, mxUpstreamAnchor(t), false},
		{"GCP Milan VCEK, the host's table (VCEK+ASK+ARK, DER)", mxToken(t, mxGCP, nil), EvidenceMediaTypeTSMCbor, gcpData, mxAMDAnchor(t, "Milan"), false},
		{"GCP Milan VCEK, a host's table with the VCEK alone", mxToken(t, mxGCP, mxLeafOnly(t, gcpTSM.AuxBlob)), EvidenceMediaTypeTSMCbor, gcpData, mxAMDAnchor(t, "Milan"), false},
		{"AWS Milan VLEK, the host's table (VLEK alone)", mxToken(t, mxAWS, mxLeafOnly(t, awsTSM.AuxBlob)), EvidenceMediaTypeTSMCbor, awsData, mxAMDAnchor(t, "Milan"), false},
		{"AWS Milan VLEK, the guest appended ASVK+ARK", mxToken(t, mxAWS, nil), EvidenceMediaTypeTSMCbor, awsData, mxAMDAnchor(t, "Milan"), false},
		{"GCP VCEK alone, then AWS VLEK alone in the same process: the VCEK", mxToken(t, mxGCP, mxLeafOnly(t, gcpTSM.AuxBlob)), EvidenceMediaTypeTSMCbor, gcpData, mxAMDAnchor(t, "Milan"), false},
		{"GCP VCEK alone, then AWS VLEK alone in the same process: the VLEK", mxToken(t, mxAWS, mxLeafOnly(t, awsTSM.AuxBlob)), EvidenceMediaTypeTSMCbor, awsData, mxAMDAnchor(t, "Milan"), true},
	}

	impl := NewImplementation()
	for _, r := range rows {
		if !r.keepCache {
			trust.ClearProductCertCache()
		}
		ct.urls = nil
		ev := &appraisal.Evidence{Data: r.data, MediaType: r.mediaType, Nonce: r.nonce}
		ids, err1 := impl.GetTrustAnchorIDs(ev)
		taID := mxShort(err1)
		if err1 == nil && len(ids) == 1 && ids[0].Class != nil && ids[0].Class.Model != nil {
			taID = "ok: " + *ids[0].Class.Model
		}
		err2 := impl.ValidateEvidenceIntegrity(ev, r.ta, nil)
		var e2 error = err2
		if err2 != nil && errors.Is(err2, ErrTAMismatch) {
			e2 = ErrTAMismatch
		}
		fmt.Printf("MX|%s|%s|%s|%s|%s|%d\n", variant, gsg, r.name, taID, mxShort(e2), len(ct.urls))
	}
}
