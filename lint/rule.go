// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lint

import (
	"regexp"
	"strings"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
)

// ProtoRule defines a lint rule that checks Google Protobuf APIs.
//
// Anything that satisfies this interface can be used as a rule,
// but most rule authors will want to use the implementations provided.
type ProtoRule interface {
	// GetName returns the name of the rule.
	GetName() RuleName

	// Lint accepts a Descriptor and lints it,
	// returning a slice of Problem objects it finds.
	Lint(desc.Descriptor) []Problem
}

// FileRule defines a lint rule that checks a file as a whole.
type FileRule struct {
	Name RuleName

	// LintFile accepts a FileDescriptor and lints it, returning a slice of
	// Problems it finds.
	LintFile func(*desc.FileDescriptor) []Problem

	// OnlyIf accepts a FileDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.FileDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *FileRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintFile on it
// only if the descriptor is a FileDescriptor.
func (r *FileRule) Lint(d desc.Descriptor) []Problem {
	if f, ok := d.(*desc.FileDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(f) {
			return r.LintFile(f)
		}
	}
	return nil
}

// MessageRule defines a lint rule that is run on each message in the file.
//
// Both top-level messages and nested messages are visited.
type MessageRule struct {
	Name RuleName

	// LintMessage accepts a MessageDescriptor and lints it, returning a slice
	// of Problems it finds.
	LintMessage func(*desc.MessageDescriptor) []Problem

	// OnlyIf accepts a MessageDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.MessageDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *MessageRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintMessage on it
// only if the descriptor is a MessageDescriptor.
func (r *MessageRule) Lint(d desc.Descriptor) []Problem {
	if m, ok := d.(*desc.MessageDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(m) {
			return r.LintMessage(m)
		}
	}
	return nil
}

// FieldRule defines a lint rule that is run on each field within a file.
type FieldRule struct {
	Name RuleName

	// LintField accepts a FieldDescriptor and lints it, returning a slice of
	// Problems it finds.
	LintField func(*desc.FieldDescriptor) []Problem

	// OnlyIf accepts a FieldDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.FieldDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *FieldRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintField on it
// only if the descriptor is a FieldDescriptor.
func (r *FieldRule) Lint(d desc.Descriptor) []Problem {
	if f, ok := d.(*desc.FieldDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(f) {
			return r.LintField(f)
		}
	}
	return nil
}

// ServiceRule defines a lint rule that is run on each service.
type ServiceRule struct {
	Name RuleName

	// LintService accepts a ServiceDescriptor and lints it.
	LintService func(*desc.ServiceDescriptor) []Problem

	// OnlyIf accepts a ServiceDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.ServiceDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *ServiceRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintSservice on it
// only if the descriptor is a ServiceDescriptor.
func (r *ServiceRule) Lint(d desc.Descriptor) []Problem {
	if s, ok := d.(*desc.ServiceDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(s) {
			return r.LintService(s)
		}
	}
	return nil
}

// MethodRule defines a lint rule that is run on each method.
type MethodRule struct {
	Name RuleName

	// LintMethod accepts a MethodDescriptor and lints it.
	LintMethod func(*desc.MethodDescriptor) []Problem

	// OnlyIf accepts a MethodDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.MethodDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *MethodRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintMethod on it
// only if the descriptor is a MethodDescriptor.
func (r *MethodRule) Lint(d desc.Descriptor) []Problem {
	if m, ok := d.(*desc.MethodDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(m) {
			return r.LintMethod(m)
		}
	}
	return nil
}

// EnumRule defines a lint rule that is run on each enum.
type EnumRule struct {
	Name RuleName

	// LintEnum accepts a EnumDescriptor and lints it.
	LintEnum func(*desc.EnumDescriptor) []Problem

	// OnlyIf accepts an EnumDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.EnumDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *EnumRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintEnum on it
// only if the descriptor is a EnumDescriptor.
func (r *EnumRule) Lint(d desc.Descriptor) []Problem {
	if e, ok := d.(*desc.EnumDescriptor); ok {
		if r.OnlyIf == nil || r.OnlyIf(e) {
			return r.LintEnum(e)
		}
	}
	return nil
}

// EnumValueRule defines a lint rule that is run on each enum value.
type EnumValueRule struct {
	Name RuleName

	// LintEnumValue accepts a EnumValueDescriptor and lints it.
	LintEnumValue func(*desc.EnumValueDescriptor) []Problem

	// OnlyIf accepts an EnumValueDescriptor and determines whether this rule
	// is applicable.
	OnlyIf func(*desc.EnumValueDescriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *EnumValueRule) GetName() RuleName {
	return r.Name
}

// Lint visits every enum value in the file and runs `LintEnum`.
//
// If an `OnlyIf` function is provided on the rule, it is run against each
// enum value, and if it returns false, the `LintEnum` function is not called.
func (r *EnumValueRule) Lint(fd *desc.FileDescriptor) []Problem {
	problems := []Problem{}

	// Lint all enums, either at the top of the file, or nested within messages.
	for _, enum := range getAllEnums(fd) {
		for _, value := range enum.GetValues() {
			if r.OnlyIf == nil || r.OnlyIf(value) {
				problems = append(problems, r.LintEnumValue(value)...)
			}
		}
	}
	return problems
}

// DescriptorRule defines a lint rule that is run on every descriptor
// in the file (but not the file itself).
type DescriptorRule struct {
	Name RuleName

	// LintDescriptor accepts a generic descriptor and lints it.
	//
	// Note: Unless the descriptor is typecast to a more specific type,
	// only a subset of methods are available to it.
	LintDescriptor func(desc.Descriptor) []Problem

	// OnlyIf accepts a Descriptor and determines whether this rule
	// is applicable.
	OnlyIf func(desc.Descriptor) bool

	noPositional struct{}
}

// GetName returns the name of the rule.
func (r *DescriptorRule) GetName() RuleName {
	return r.Name
}

// Lint accepts a Descriptor and applies LintDescriptor on it.
func (r *DescriptorRule) Lint(d desc.Descriptor) []Problem {
	if r.OnlyIf == nil || r.OnlyIf(d) {
		return r.LintDescriptor(d)
	}
	return nil
}

var disableRuleNameRegex = regexp.MustCompile(`api-linter:\s*(.+)\s*=\s*disabled`)

func extractDisabledRuleName(commentLine string) string {
	match := disableRuleNameRegex.FindStringSubmatch(commentLine)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

// ruleIsEnabled returns true if the rule is enabled (not disabled by the comments
// for the given descriptor or its file), false otherwise.
func ruleIsEnabled(rule ProtoRule, d desc.Descriptor, aliasMap map[string]string) bool {
	// Some rules have a legacy name. We add it to the check list.
	ruleName := string(rule.GetName())
	names := []string{ruleName, aliasMap[ruleName]}

	commentLines := strings.Split(fileHeader(d.GetFile()), "\n")
	commentLines = append(commentLines, strings.Split(getLeadingComments(d), "\n")...)
	disabledRules := []string{}
	for _, commentLine := range commentLines {
		r := extractDisabledRuleName(commentLine)
		if r != "" {
			disabledRules = append(disabledRules, r)
		}
	}

	for _, name := range names {
		if matchRule(name, disabledRules...) {
			return false
		}
	}

	return true
}

func getLeadingComments(d desc.Descriptor) string {
	if sourceInfo := d.GetSourceInfo(); sourceInfo != nil {
		return sourceInfo.GetLeadingComments()
	}
	return ""
}

// getAllDescriptors returns a slice with every descriptor in the file.
func getAllDescriptors(f *desc.FileDescriptor) []desc.Descriptor {
	descriptors := []desc.Descriptor{f}
	for _, d := range getAllEnums(f) {
		descriptors = append(descriptors, d)
	}
	for _, d := range getAllEnumValues(f) {
		descriptors = append(descriptors, d)
	}
	for _, d := range getAllFields(f) {
		descriptors = append(descriptors, d)
	}
	for _, d := range getAllMethods(f) {
		descriptors = append(descriptors, d)
	}
	for _, d := range getAllMessages(f) {
		descriptors = append(descriptors, d)
	}
	for _, d := range f.GetServices() {
		descriptors = append(descriptors, d)
	}
	return descriptors
}

// getAllMethods returns a slice with every method in the file.
func getAllMethods(f *desc.FileDescriptor) []*desc.MethodDescriptor {
	methods := []*desc.MethodDescriptor{}
	for _, service := range f.GetServices() {
		methods = append(methods, service.GetMethods()...)
	}
	return methods
}

// getAllMessages returns a slice with every message (not just top-level
// messages) in the file.
func getAllMessages(f *desc.FileDescriptor) (messages []*desc.MessageDescriptor) {
	messages = append(messages, f.GetMessageTypes()...)
	for _, message := range f.GetMessageTypes() {
		messages = append(messages, getAllNestedMessages(message)...)
	}
	return messages
}

// getAllNestedMessages returns a slice with the given message descriptor as well
// as all nested message descriptors, traversing to arbitrary depth.
func getAllNestedMessages(m *desc.MessageDescriptor) (messages []*desc.MessageDescriptor) {
	for _, nested := range m.GetNestedMessageTypes() {
		if !nested.IsMapEntry() { // Don't include the synthetic message type that represents an entry in a map field.
			messages = append(messages, nested)
		}
		messages = append(messages, getAllNestedMessages(nested)...)
	}
	return messages
}

// getAllEnums returns a slice with every enum (not just top-level enums)
// in the file.
func getAllEnums(f *desc.FileDescriptor) (enums []*desc.EnumDescriptor) {
	// Append all enums at the top level.
	enums = append(enums, f.GetEnumTypes()...)

	// Append all enums nested within messages.
	for _, m := range getAllMessages(f) {
		enums = append(enums, m.GetNestedEnumTypes()...)
	}

	return
}

// getAllEnumValues returns a slice with every enum value in the file.
func getAllEnumValues(f *desc.FileDescriptor) []*desc.EnumValueDescriptor {
	values := []*desc.EnumValueDescriptor{}
	for _, e := range getAllEnums(f) {
		values = append(values, e.GetValues()...)
	}
	return values
}

// getAllFields returns a slice with every field in the file.
func getAllFields(f *desc.FileDescriptor) []*desc.FieldDescriptor {
	fields := []*desc.FieldDescriptor{}
	for _, m := range getAllMessages(f) {
		fields = append(fields, m.GetFields()...)
	}
	return fields
}

// fileHeader attempts to get the comment at the top of the file, but it
// is on a best effort basis because protobuf is inconsistent.
//
// Taken from https://github.com/jhump/protoreflect/issues/215
func fileHeader(fd *desc.FileDescriptor) string {
	var firstLoc *dpb.SourceCodeInfo_Location
	var firstSpan int64

	// File level comments should only be comments identified on either
	// syntax (12), package (2), option (8), or import (3) statements.
	allowedPaths := map[int32]struct{}{2: {}, 3: {}, 8: {}, 12: {}}

	// Iterate over locations in the file descriptor looking for
	// what we think is a file-level comment.
	for _, curr := range fd.AsFileDescriptorProto().GetSourceCodeInfo().GetLocation() {
		// Skip locations that have no comments.
		if curr.LeadingComments == nil && len(curr.LeadingDetachedComments) == 0 {
			continue
		}
		// Skip locations that are not allowed because they should never be
		// mistaken for file-level comments.
		if _, ok := allowedPaths[curr.GetPath()[0]]; !ok {
			continue
		}
		currSpan := asPos(curr.Span)
		if firstLoc == nil || currSpan < firstSpan {
			firstLoc = curr
			firstSpan = currSpan
		}
	}
	if firstLoc == nil {
		return ""
	}
	if len(firstLoc.LeadingDetachedComments) > 0 {
		return strings.Join(firstLoc.LeadingDetachedComments, "\n")
	}
	return firstLoc.GetLeadingComments()
}

func asPos(span []int32) int64 {
	return (int64(span[0]) << 32) + int64(span[1])
}
