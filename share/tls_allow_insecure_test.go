package share

import (
	"testing"

	"github.com/xtls/xray-core/infra/conf"
)

func TestConvertShareLinksToXrayJsonPreservesTLSAllowInsecure(t *testing.T) {
	link := "vless://11111111-1111-1111-1111-111111111111@example.com:443?type=ws&security=tls&allowInsecure=1#tls-node"

	xray, err := ConvertShareLinksToXrayJson(link)
	if err != nil {
		t.Fatalf("ConvertShareLinksToXrayJson() error = %v", err)
	}
	if len(xray.OutboundConfigs) != 1 {
		t.Fatalf("expected 1 outbound, got %d", len(xray.OutboundConfigs))
	}

	streamSettings := xray.OutboundConfigs[0].StreamSetting
	if streamSettings == nil || streamSettings.TLSSettings == nil {
		t.Fatalf("expected TLS stream settings, got %#v", streamSettings)
	}
	if !streamSettings.TLSSettings.AllowInsecure {
		t.Fatal("expected TLS AllowInsecure to be true")
	}
}

func TestShareLinkIncludesTLSAllowInsecure(t *testing.T) {
	settingsRawMessage, err := convertJsonToRawMessage(conf.VLessOutboundConfig{
		Address:    parseAddress("example.com"),
		Port:       443,
		Id:         "11111111-1111-1111-1111-111111111111",
		Encryption: "none",
	})
	if err != nil {
		t.Fatalf("convertJsonToRawMessage() error = %v", err)
	}
	network := conf.TransportProtocol("ws")

	link, err := shareLink(conf.OutboundDetourConfig{
		Protocol: "vless",
		Settings: &settingsRawMessage,
		StreamSetting: &conf.StreamConfig{
			Network:  &network,
			Security: "tls",
			TLSSettings: &conf.TLSConfig{
				AllowInsecure: true,
			},
		},
	})
	if err != nil {
		t.Fatalf("shareLink() error = %v", err)
	}

	if got := link.Query().Get("allowInsecure"); got != "1" {
		t.Fatalf("allowInsecure query = %q, want 1", got)
	}
}
