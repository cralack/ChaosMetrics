// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: publisher.proto

package publisher

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TaskSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Loc          string `protobuf:"bytes,3,opt,name=loc,proto3" json:"loc,omitempty"`
	AssignedNode string `protobuf:"bytes,4,opt,name=assigned_node,json=assignedNode,proto3" json:"assigned_node,omitempty"`
	CreationTime int64  `protobuf:"varint,5,opt,name=creation_time,json=creationTime,proto3" json:"creation_time,omitempty"`
	Sumname      string `protobuf:"bytes,6,opt,name=sumname,proto3" json:"sumname,omitempty"`
	Type         string `protobuf:"bytes,7,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *TaskSpec) Reset() {
	*x = TaskSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_publisher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSpec) ProtoMessage() {}

func (x *TaskSpec) ProtoReflect() protoreflect.Message {
	mi := &file_publisher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSpec.ProtoReflect.Descriptor instead.
func (*TaskSpec) Descriptor() ([]byte, []int) {
	return file_publisher_proto_rawDescGZIP(), []int{0}
}

func (x *TaskSpec) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TaskSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TaskSpec) GetLoc() string {
	if x != nil {
		return x.Loc
	}
	return ""
}

func (x *TaskSpec) GetAssignedNode() string {
	if x != nil {
		return x.AssignedNode
	}
	return ""
}

func (x *TaskSpec) GetCreationTime() int64 {
	if x != nil {
		return x.CreationTime
	}
	return 0
}

func (x *TaskSpec) GetSumname() string {
	if x != nil {
		return x.Sumname
	}
	return ""
}

func (x *TaskSpec) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

var File_publisher_proto protoreflect.FileDescriptor

var file_publisher_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb8, 0x01, 0x0a,
	0x08, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6c, 0x6f, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x6f, 0x63, 0x12,
	0x23, 0x0a, 0x0d, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x6e, 0x6f, 0x64, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x32, 0x56, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x73, 0x68, 0x65, 0x72, 0x12, 0x49, 0x0a, 0x08, 0x50, 0x75, 0x73, 0x68, 0x54, 0x61, 0x73, 0x6b,
	0x12, 0x09, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x70, 0x65, 0x63, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x22, 0x0f,
	0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x42,
	0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_publisher_proto_rawDescOnce sync.Once
	file_publisher_proto_rawDescData = file_publisher_proto_rawDesc
)

func file_publisher_proto_rawDescGZIP() []byte {
	file_publisher_proto_rawDescOnce.Do(func() {
		file_publisher_proto_rawDescData = protoimpl.X.CompressGZIP(file_publisher_proto_rawDescData)
	})
	return file_publisher_proto_rawDescData
}

var file_publisher_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_publisher_proto_goTypes = []interface{}{
	(*TaskSpec)(nil),    // 0: TaskSpec
	(*empty.Empty)(nil), // 1: google.protobuf.Empty
}
var file_publisher_proto_depIdxs = []int32{
	0, // 0: Publisher.PushTask:input_type -> TaskSpec
	1, // 1: Publisher.PushTask:output_type -> google.protobuf.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_publisher_proto_init() }
func file_publisher_proto_init() {
	if File_publisher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_publisher_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskSpec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_publisher_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_publisher_proto_goTypes,
		DependencyIndexes: file_publisher_proto_depIdxs,
		MessageInfos:      file_publisher_proto_msgTypes,
	}.Build()
	File_publisher_proto = out.File
	file_publisher_proto_rawDesc = nil
	file_publisher_proto_goTypes = nil
	file_publisher_proto_depIdxs = nil
}
