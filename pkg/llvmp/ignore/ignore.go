package ignore

import (
	"fmt"
	"regexp"
	"strings"
)

type Set map[string]bool

func (m Set) Match(fname string) bool {
	if _, ok := m[fname]; ok {
		return true
	}
	return false
}

var (
	ignoredFn = map[string][]string{
		"@default": {
			"@bpf",
			"@builtins",
			"@cildbg",
			"@ctx",
			"@eth",
			"@ipv4",
			"@lb",
			"@metrics",
			"@srv6",
			"@utils",
		},
		"@bpf": {
			"__bpf_memcpy",
			"bpf_barrier",
			"bpf_clear_meta",
		},
		"@builtins": {
			"llvm",
			"memcpy",
			"memset",
		},
		"@cildbg": {
			"cilium_dbg",
			"cilium_dbg3",
			"cilium_dbg_capture",
			"cilium_capture_out",
		},
		"@ctx": {
			"ctx_change_head",
			"ctx_data_end",
			"ctx_data_start",
			"ctx_data",
			"ctx_full_len",
			"ctx_get_ifindex",
			"ctx_get_protocol",
			"ctx_is_skb",
			"ctx_load_and_clear_meta",
			"ctx_load_meta",
			"ctx_redirect",
			"ctx_set_encap_info",
			"ctx_store_meta",
		},
		"@eth": {
			"eth_addrcmp",
			"eth_is_bcast",
			"eth_is_supported_ethertype",
			"eth_load_saddr",
			"eth_store_saddr_aligned",
			"eth_store_saddr",
			"eth_load_daddr",
			"eth_store_daddr_aligned",
			"eth_store_daddr",
			"eth_store_proto",
		},
		"@lb": {
			"lb4_svc_is_affinity",
			"lb4_update_affinity_by_addr",
			"lb4_fill_key",
			"lb4_lookup_service",
			"lb4_affinity_backend_id_by_addr",
			"lb4_lookup_backend",
		},
		"@ipv4": {
			"ipv4_has_l4_header",
			"ipv4_hdrlen",
			"ipv4_load_l4_ports",
			"ipv4_load_l4_ports",
		},
		"@metrics": {
			"_send_trace_notify",
			"_update_metrics",
		},
		"@utils": {
			"__id_for_file",
			"__revalidate_data_pull",
			"csum_diff",
			"csum_l4_replace",
			"ipv6_addr_copy",
			"_utime_get_offset",
			"utime_get_time",
			"is_valid_lxc_src_ipv4",
		},
		"@srv6": {
			"srv6_lookup_vrf4",
			"srv6_lookup_policy4",
			"srv6_lookup_vrf6",
			"srv6_lookup_policy6",
			"srv6_lookup_sid",
			"srv6_encapsulation",
			"srv6_decapsulation",
			"srv6_handling4",
			"srv6_handling6",
			"srv6_handling",
			"srv6_load_meta_sid",
			"srv6_store_meta_sid",
		},
	}
	ignoredFnRe = map[string][]regexp.Regexp{}
)

func expand(l []string) ([]string, error) {
	var ret []string
	for _, x := range l {
		if strings.HasPrefix(x, "@") {
			entry, ok := ignoredFn[x]
			if !ok {
				return nil, fmt.Errorf("fn list %q does not exist", x)
			}
			// TODO: need to error out on infinite loop.
			el, err := expand(entry)
			if err != nil {
				return nil, err
			}
			ret = append(ret, el...)
		} else {
			ret = append(ret, x)
		}
	}
	return ret, nil
}

func Make(l []string) (Set, error) {
	m := Set{}
	el, err := expand(l)
	if err != nil {
		return nil, err
	}
	for _, x := range el {
		m[x] = true
	}
	return m, nil
}
