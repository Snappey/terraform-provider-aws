// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=GetTags -ServiceTagsMap -TagInTagsElem=TagsToAdd -UntagInTagsElem=TagsToRemove -UpdateTags
// ONLY generate directives and package declaration! Do not add anything else to this file.

package glue
