package serializer

import (
	"testing"

	"github.com/hemanthkumarkola1/gRPCProj/pb"
	"github.com/hemanthkumarkola1/gRPCProj/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"

	latop1 := sample.NewLaptop()
	err := WriteProtobufToBinaryFile(latop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)

	require.NoError(t, err)

	require.True(t, proto.Equal(laptop2, latop1))

	WriteProtobufToBinaryFile(nil, binaryFile) // To clear the written binary file.
}
