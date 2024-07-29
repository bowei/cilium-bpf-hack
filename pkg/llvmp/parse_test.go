package llvmp

import (
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	for _, tc := range []struct {
		name       string
		re         *regexp.Regexp
		matches    []string
		notMatches []string
	}{
		{
			name: "fnStartRe",
			re:   fnStartRe,
			matches: []string{
				`define dso_local i32 @__send_drop_notify(ptr noundef %0) #0 section "2/1" !dbg !2036 {`,
			},
		},
		{
			name: "fnStart2Re",
			re:   fnStart2Re,
			matches: []string{
				`define dso_local i32 @__send_drop_notify(ptr noundef %0) #0 section "2/1" !dbg !2036 {`,
			},
		},
		{
			name: "fnSectionRe",
			re:   fnSectionRe,
			matches: []string{
				` section "2/1" `,
			},
		},
		{
			name: "fnEndRe",
			re:   fnEndRe,
			matches: []string{
				`}`,
			},
			notMatches: []string{
				` abc}`,
			},
		},
		{
			name: "callRe",
			re:   callRe,
			matches: []string{
				`  %12 = call i32 @srv6_decapsulation(ptr noundef %11), !dbg !2907`,
				`  call void @llvm.dbg.declare(metadata ptr %3, metadata !2902, metadata !DIExpression()), !dbg !2903`,
				`  %38 = call ptr @ctx_data_end(ptr noundef %37), !dbg !3365`,
			},
		},
		{
			name: "callIndirectPtrRe",
			re:   callIndirectRe,
			matches: []string{
				`  %21 = call ptr %18(ptr noundef %19, ptr noundef %20), !dbg !16832`,
			},
		},
		{
			name: "callSymRe",
			re:   callSymRe,
			matches: []string{
				`  %38 = call ptr @ctx_data_end(ptr noundef %37), !dbg !3365`,
				`  call void @ipv6_addr_copy(ptr noundef %145, ptr noundef %146), !dbg !3414`,
			},
		},
		{
			name: "tcInternalRe",
			re:   tcInternalRe,
			matches: []string{
				`  %63 = call i32 @tail_call_internal(ptr noundef %62, i32 noundef 7, ptr noundef %7), !dbg !5886`,
			},
		},
		{
			name: "tcDyanmicRe",
			re:   tcDyanmicRe,
			matches: []string{
				`  call void @tail_call_dynamic(ptr noundef %5, ptr noundef @POLICY_EGRESSCALL_MAP, i32 noundef %7), !dbg !10557`,
			},
		},
		{
			name: "tcPolicyRe",
			re:   tcPolicyRe,
			matches: []string{
				`  %46 = call i32 @tail_call_policy(ptr noundef %42, i16 noundef zeroext %45), !dbg !17861`,
			},
		},
		{
			name: "tcEgressPolicyRe",
			re:   tcEgressPolicyRe,
			matches: []string{
				`  %53 = call i32 @tail_call_egress_policy(ptr noundef %50, i16 noundef zeroext %52), !dbg !10370`,
			},
		},
		{
			name: "diLocationRe",
			re:   diLocationRe,
			matches: []string{
				`!21084 = !DILocation(line: 76, column: 34, scope: !21082)`,
			},
		},
		{
			name: "diLexicalBlockRe",
			re:   diLexicalBlockRe,
			matches: []string{
				`!18333 = distinct !DILexicalBlock(scope: !18334, file: !3, line: 169, column: 7)`,
			},
		},
		{
			name: "diSubprogramRe",
			re:   diSubprogramRe,
			matches: []string{
				`!14753 = distinct !DISubprogram(name: "ct_has_nodeport_egress_entry6", scope: !227, file: !227, line: 1140, type: !14754, scopeLine: 1143, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition, unit: !2, retainedNodes: !2040)`,
			},
		},
		{
			name: "diFileRe",
			re:   diFileRe,
			matches: []string{
				`!21090 = !DIFile(filename: "lib/clustermesh.h", directory: "/home/bowei/work/cilium/bpf", checksumkind: CSK_MD5, checksum: "2a78f5e0ef473d44b386802cd1dcefd7")`,
			},
		},
		{
			name: "diSubprogramRe",
			re:   diSubprogramRe,
			// TODO
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, str := range tc.matches {
				if !tc.re.Match([]byte(str)) {
					t.Errorf("Match(%q) = false, want true", str)
				}
			}
			for _, str := range tc.notMatches {
				if tc.re.Match([]byte(str)) {
					t.Errorf("Match(%q) = true, want false", str)
				}
			}
		})
	}
}
