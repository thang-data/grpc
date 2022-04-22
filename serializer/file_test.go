package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/dem1/dem1/pb"
	"gitlab.com/dem1/dem1/sample"
	"gitlab.com/dem1/dem1/serializer"
)

func TestSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()

	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	err = serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
}
