<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [proofs.proto](#proofs.proto)
    - [BatchEntry](#ics23.BatchEntry)
    - [BatchProof](#ics23.BatchProof)
    - [CommitmentProof](#ics23.CommitmentProof)
    - [CompressedBatchEntry](#ics23.CompressedBatchEntry)
    - [CompressedBatchProof](#ics23.CompressedBatchProof)
    - [CompressedExistenceProof](#ics23.CompressedExistenceProof)
    - [CompressedNonExistenceProof](#ics23.CompressedNonExistenceProof)
    - [ExistenceProof](#ics23.ExistenceProof)
    - [InnerOp](#ics23.InnerOp)
    - [InnerSpec](#ics23.InnerSpec)
    - [LeafOp](#ics23.LeafOp)
    - [NonExistenceProof](#ics23.NonExistenceProof)
    - [ProofSpec](#ics23.ProofSpec)
  
    - [HashOp](#ics23.HashOp)
    - [LengthOp](#ics23.LengthOp)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proofs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proofs.proto



<a name="ics23.BatchEntry"></a>

### BatchEntry
Use BatchEntry not CommitmentProof, to avoid recursion


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exist` | [ExistenceProof](#ics23.ExistenceProof) |  |  |
| `nonexist` | [NonExistenceProof](#ics23.NonExistenceProof) |  |  |






<a name="ics23.BatchProof"></a>

### BatchProof
BatchProof is a group of multiple proof types than can be compressed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [BatchEntry](#ics23.BatchEntry) | repeated |  |






<a name="ics23.CommitmentProof"></a>

### CommitmentProof
CommitmentProof is either an ExistenceProof or a NonExistenceProof, or a Batch of such messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exist` | [ExistenceProof](#ics23.ExistenceProof) |  |  |
| `nonexist` | [NonExistenceProof](#ics23.NonExistenceProof) |  |  |
| `batch` | [BatchProof](#ics23.BatchProof) |  |  |
| `compressed` | [CompressedBatchProof](#ics23.CompressedBatchProof) |  |  |






<a name="ics23.CompressedBatchEntry"></a>

### CompressedBatchEntry
Use BatchEntry not CommitmentProof, to avoid recursion


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exist` | [CompressedExistenceProof](#ics23.CompressedExistenceProof) |  |  |
| `nonexist` | [CompressedNonExistenceProof](#ics23.CompressedNonExistenceProof) |  |  |






<a name="ics23.CompressedBatchProof"></a>

### CompressedBatchProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `entries` | [CompressedBatchEntry](#ics23.CompressedBatchEntry) | repeated |  |
| `lookup_inners` | [InnerOp](#ics23.InnerOp) | repeated |  |






<a name="ics23.CompressedExistenceProof"></a>

### CompressedExistenceProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |
| `leaf` | [LeafOp](#ics23.LeafOp) |  |  |
| `path` | [int32](#int32) | repeated | these are indexes into the lookup_inners table in CompressedBatchProof |






<a name="ics23.CompressedNonExistenceProof"></a>

### CompressedNonExistenceProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | TODO: remove this as unnecessary??? we prove a range |
| `left` | [CompressedExistenceProof](#ics23.CompressedExistenceProof) |  |  |
| `right` | [CompressedExistenceProof](#ics23.CompressedExistenceProof) |  |  |






<a name="ics23.ExistenceProof"></a>

### ExistenceProof
ExistenceProof takes a key and a value and a set of steps to perform on it.
The result of peforming all these steps will provide a "root hash", which can
be compared to the value in a header.
Since it is computationally infeasible to produce a hash collission for any of the used
cryptographic hash functions, if someone can provide a series of operations to transform
a given key and value into a root hash that matches some trusted root, these key and values
must be in the referenced merkle tree.
The only possible issue is maliablity in LeafOp, such as providing extra prefix data,
which should be controlled by a spec. Eg. with lengthOp as NONE,
prefix = FOO, key = BAR, value = CHOICE
and
prefix = F, key = OOBAR, value = CHOICE
would produce the same value.
With LengthOp this is tricker but not impossible. Which is why the "leafPrefixEqual" field
in the ProofSpec is valuable to prevent this mutability. And why all trees should
length-prefix the data before hashing it.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |
| `leaf` | [LeafOp](#ics23.LeafOp) |  |  |
| `path` | [InnerOp](#ics23.InnerOp) | repeated |  |






<a name="ics23.InnerOp"></a>

### InnerOp
InnerOp represents a merkle-proof step that is not a leaf.
It represents concatenating two children and hashing them to provide the next result.
The result of the previous step is passed in, so the signature of this op is:
innerOp(child) -> output
The result of applying InnerOp should be:
output = op.hash(op.prefix || child || op.suffix)
where the || operator is concatenation of binary data,
and child is the result of hashing all the tree below this step.
Any special data, like prepending child with the length, or prepending the entire operation with
some value to differentiate from leaf nodes, should be included in prefix and suffix.
If either of prefix or suffix is empty, we just treat it as an empty string


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [HashOp](#ics23.HashOp) |  |  |
| `prefix` | [bytes](#bytes) |  |  |
| `suffix` | [bytes](#bytes) |  |  |






<a name="ics23.InnerSpec"></a>

### InnerSpec
InnerSpec contains all store-specific structure info to determine if two proofs from a
given store are neighbors.
This enables:
isLeftMost(spec: InnerSpec, op: InnerOp)
isRightMost(spec: InnerSpec, op: InnerOp)
isLeftNeighbor(spec: InnerSpec, left: InnerOp, right: InnerOp)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `child_order` | [int32](#int32) | repeated | Child order is the ordering of the children node, must count from 0 iavl tree is [0, 1] (left then right) merk is [0, 2, 1] (left, right, here) |
| `child_size` | [int32](#int32) |  |  |
| `min_prefix_length` | [int32](#int32) |  |  |
| `max_prefix_length` | [int32](#int32) |  |  |
| `empty_child` | [bytes](#bytes) |  | empty child is the prehash image that is used when one child is nil (eg. 20 bytes of 0) |
| `hash` | [HashOp](#ics23.HashOp) |  | hash is the algorithm that must be used for each InnerOp |






<a name="ics23.LeafOp"></a>

### LeafOp
LeafOp represents the raw key-value data we wish to prove, and
must be flexible to represent the internal transformation from
the original key-value pairs into the basis hash, for many existing
merkle trees.
key and value are passed in. So that the signature of this operation is:
leafOp(key, value) -> output
To process this, first prehash the keys and values if needed (ANY means no hash in this case):
hkey = prehashKey(key)
hvalue = prehashValue(value)
Then combine the bytes, and hash it
output = hash(prefix || length(hkey) || hkey || length(hvalue) || hvalue)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [HashOp](#ics23.HashOp) |  |  |
| `prehash_key` | [HashOp](#ics23.HashOp) |  |  |
| `prehash_value` | [HashOp](#ics23.HashOp) |  |  |
| `length` | [LengthOp](#ics23.LengthOp) |  |  |
| `prefix` | [bytes](#bytes) |  | prefix is a fixed bytes that may optionally be included at the beginning to differentiate a leaf node from an inner node. |






<a name="ics23.NonExistenceProof"></a>

### NonExistenceProof
NonExistenceProof takes a proof of two neighbors, one left of the desired key,
one right of the desired key. If both proofs are valid AND they are neighbors,
then there is no valid proof for the given key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | TODO: remove this as unnecessary??? we prove a range |
| `left` | [ExistenceProof](#ics23.ExistenceProof) |  |  |
| `right` | [ExistenceProof](#ics23.ExistenceProof) |  |  |






<a name="ics23.ProofSpec"></a>

### ProofSpec
ProofSpec defines what the expected parameters are for a given proof type.
This can be stored in the client and used to validate any incoming proofs.
verify(ProofSpec, Proof) -> Proof | Error
As demonstrated in tests, if we don't fix the algorithm used to calculate the
LeafHash for a given tree, there are many possible key-value pairs that can
generate a given hash (by interpretting the preimage differently).
We need this for proper security, requires client knows a priori what
tree format server uses. But not in code, rather a configuration object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `leaf_spec` | [LeafOp](#ics23.LeafOp) |  | any field in the ExistenceProof must be the same as in this spec. except Prefix, which is just the first bytes of prefix (spec can be longer) |
| `inner_spec` | [InnerSpec](#ics23.InnerSpec) |  |  |
| `max_depth` | [int32](#int32) |  | max_depth (if > 0) is the maximum number of InnerOps allowed (mainly for fixed-depth tries) |
| `min_depth` | [int32](#int32) |  | min_depth (if > 0) is the minimum number of InnerOps allowed (mainly for fixed-depth tries) |





 <!-- end messages -->


<a name="ics23.HashOp"></a>

### HashOp


| Name | Number | Description |
| ---- | ------ | ----------- |
| NO_HASH | 0 | NO_HASH is the default if no data passed. Note this is an illegal argument some places. |
| SHA256 | 1 |  |
| SHA512 | 2 |  |
| KECCAK | 3 |  |
| RIPEMD160 | 4 |  |
| BITCOIN | 5 | ripemd160(sha256(x)) |
| SHA512_256 | 6 |  |



<a name="ics23.LengthOp"></a>

### LengthOp
LengthOp defines how to process the key and value of the LeafOp
to include length information. After encoding the length with the given
algorithm, the length will be prepended to the key and value bytes.
(Each one with it's own encoded length)

| Name | Number | Description |
| ---- | ------ | ----------- |
| NO_PREFIX | 0 | NO_PREFIX don't include any length info |
| VAR_PROTO | 1 | VAR_PROTO uses protobuf (and go-amino) varint encoding of the length |
| VAR_RLP | 2 | VAR_RLP uses rlp int encoding of the length |
| FIXED32_BIG | 3 | FIXED32_BIG uses big-endian encoding of the length as a 32 bit integer |
| FIXED32_LITTLE | 4 | FIXED32_LITTLE uses little-endian encoding of the length as a 32 bit integer |
| FIXED64_BIG | 5 | FIXED64_BIG uses big-endian encoding of the length as a 64 bit integer |
| FIXED64_LITTLE | 6 | FIXED64_LITTLE uses little-endian encoding of the length as a 64 bit integer |
| REQUIRE_32_BYTES | 7 | REQUIRE_32_BYTES is like NONE, but will fail if the input is not exactly 32 bytes (sha256 output) |
| REQUIRE_64_BYTES | 8 | REQUIRE_64_BYTES is like NONE, but will fail if the input is not exactly 64 bytes (sha512 output) |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

