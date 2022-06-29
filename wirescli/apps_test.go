package wirescli

import (
	pprint "github.com/NubeIO/rubix-assist/pkg/helpers/print"
	"testing"
)

func TestClient_InstallApp(t *testing.T) {

	client := New("0.0.0.0", 1313)

	body := &WiresTokenBody{
		Username: "admin",
		Password: "admin",
	}
	token, _ := client.GetToken(body)
	pprint.PrintJOSN(token)
	data, _ := client.Backup(token.Token)
	pprint.PrintJOSN(data.Objects)
}
