# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: recognizer.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='recognizer.proto',
  package='proto',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=_b('\n\x10recognizer.proto\x12\x05proto\"\x1f\n\x10RecognizeRequest\x12\x0b\n\x03url\x18\x01 \x01(\t\"M\n\x04\x46\x61\x63\x65\x12\t\n\x01x\x18\x01 \x01(\x05\x12\t\n\x01y\x18\x02 \x01(\x05\x12\r\n\x05width\x18\x03 \x01(\x05\x12\x0e\n\x06height\x18\x04 \x01(\x05\x12\x10\n\x08\x65ncoding\x18\x05 \x03(\x02\"/\n\x11RegognizeResponse\x12\x1a\n\x05\x66\x61\x63\x65s\x18\x01 \x03(\x0b\x32\x0b.proto.Face2U\n\x0e\x46\x61\x63\x65Recognizer\x12\x43\n\x0eRecognizeFaces\x12\x17.proto.RecognizeRequest\x1a\x18.proto.RegognizeResponseb\x06proto3')
)




_RECOGNIZEREQUEST = _descriptor.Descriptor(
  name='RecognizeRequest',
  full_name='proto.RecognizeRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='url', full_name='proto.RecognizeRequest.url', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=27,
  serialized_end=58,
)


_FACE = _descriptor.Descriptor(
  name='Face',
  full_name='proto.Face',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='x', full_name='proto.Face.x', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='y', full_name='proto.Face.y', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='width', full_name='proto.Face.width', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='height', full_name='proto.Face.height', index=3,
      number=4, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='encoding', full_name='proto.Face.encoding', index=4,
      number=5, type=2, cpp_type=6, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=60,
  serialized_end=137,
)


_REGOGNIZERESPONSE = _descriptor.Descriptor(
  name='RegognizeResponse',
  full_name='proto.RegognizeResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='faces', full_name='proto.RegognizeResponse.faces', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=139,
  serialized_end=186,
)

_REGOGNIZERESPONSE.fields_by_name['faces'].message_type = _FACE
DESCRIPTOR.message_types_by_name['RecognizeRequest'] = _RECOGNIZEREQUEST
DESCRIPTOR.message_types_by_name['Face'] = _FACE
DESCRIPTOR.message_types_by_name['RegognizeResponse'] = _REGOGNIZERESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

RecognizeRequest = _reflection.GeneratedProtocolMessageType('RecognizeRequest', (_message.Message,), dict(
  DESCRIPTOR = _RECOGNIZEREQUEST,
  __module__ = 'recognizer_pb2'
  # @@protoc_insertion_point(class_scope:proto.RecognizeRequest)
  ))
_sym_db.RegisterMessage(RecognizeRequest)

Face = _reflection.GeneratedProtocolMessageType('Face', (_message.Message,), dict(
  DESCRIPTOR = _FACE,
  __module__ = 'recognizer_pb2'
  # @@protoc_insertion_point(class_scope:proto.Face)
  ))
_sym_db.RegisterMessage(Face)

RegognizeResponse = _reflection.GeneratedProtocolMessageType('RegognizeResponse', (_message.Message,), dict(
  DESCRIPTOR = _REGOGNIZERESPONSE,
  __module__ = 'recognizer_pb2'
  # @@protoc_insertion_point(class_scope:proto.RegognizeResponse)
  ))
_sym_db.RegisterMessage(RegognizeResponse)



_FACERECOGNIZER = _descriptor.ServiceDescriptor(
  name='FaceRecognizer',
  full_name='proto.FaceRecognizer',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=188,
  serialized_end=273,
  methods=[
  _descriptor.MethodDescriptor(
    name='RecognizeFaces',
    full_name='proto.FaceRecognizer.RecognizeFaces',
    index=0,
    containing_service=None,
    input_type=_RECOGNIZEREQUEST,
    output_type=_REGOGNIZERESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_FACERECOGNIZER)

DESCRIPTOR.services_by_name['FaceRecognizer'] = _FACERECOGNIZER

# @@protoc_insertion_point(module_scope)
