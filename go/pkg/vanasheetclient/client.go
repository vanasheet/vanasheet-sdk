package vanasheetclient

import (
	"context"

	"github.com/vanasheet/vanasheet-sdk/go/pkg/vanasheetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

// VanasheetIO is the main struct that implements all vanasheet implementations
type VanasheetIO struct {
	apiKey string
	client vanasheetpb.VanasheetClient
}

// New creates a new instance of VanasheetIO
func New(apikey string, grpcConn *grpc.ClientConn) *VanasheetIO {
	return &VanasheetIO{apikey, vanasheetpb.NewVanasheetClient(grpcConn)}
}

// AppendRow appends a single row to spreadsheet and returns the row value
func (client *VanasheetIO) AppendRow(
	ctx context.Context,
	spreadsheetID string,
	sheetname string,
	appendRow map[string]interface{},
) (map[string]interface{}, error) {
	// append apikey
	outgoingCtx := metadata.AppendToOutgoingContext(ctx, "authorization", client.apiKey)

	// parse
	appendRowPb, err := structpb.NewStruct(appendRow)
	if err != nil {
		return nil, err
	}

	// send request
	pbResp, err := client.client.AppendRow(
		outgoingCtx,
		&vanasheetpb.AppendRowRequest{
			SpreadsheetId: spreadsheetID,
			Sheetname:     sheetname,
			Row:           appendRowPb,
		},
	)
	if err != nil {
		return nil, err
	}

	return pbResp.Row.AsMap(), nil
}

// RawReadQuery returns the result of a query in raw format
func (client *VanasheetIO) RawReadQuery(
	ctx context.Context,
	spreadsheetID string,
	querysheetname string,
	query string,
) ([][]interface{}, error) {
	outgoingCtx := metadata.AppendToOutgoingContext(ctx, "authorization", client.apiKey)

	respPb, err := client.client.RawReadQuery(
		outgoingCtx,
		&vanasheetpb.RawReadQueryRequest{
			SpreadsheetId: spreadsheetID,
			A1Range:       querysheetname,
			Query:         query,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := make([][]interface{}, len(respPb.Rows))

	for i, row := range respPb.Rows {
		resp[i] = make([]interface{}, len(row.Vals))
		for j, v := range row.Vals {
			resp[i][j] = v.AsInterface()
		}
	}

	return resp, nil
}
