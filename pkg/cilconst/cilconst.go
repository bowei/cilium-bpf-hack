package cilconst

// Taken from bpf/lib/common.h.
const (
	CILIUM_CALL_DROP_NOTIFY                 = 1
	CILIUM_CALL_ERROR_NOTIFY                = 2
	CILIUM_CALL_HANDLE_ICMP6_NS             = 4
	CILIUM_CALL_SEND_ICMP6_TIME_EXCEEDED    = 5
	CILIUM_CALL_ARP                         = 6
	CILIUM_CALL_IPV4_FROM_LXC               = 7
	CILIUM_CALL_IPV4_FROM_NETDEV            = CILIUM_CALL_IPV4_FROM_LXC
	CILIUM_CALL_IPV4_FROM_OVERLAY           = CILIUM_CALL_IPV4_FROM_LXC
	CILIUM_CALL_IPV46_RFC8215               = 8
	CILIUM_CALL_IPV64_RFC8215               = 9
	CILIUM_CALL_IPV6_FROM_LXC               = 10
	CILIUM_CALL_IPV6_FROM_NETDEV            = CILIUM_CALL_IPV6_FROM_LXC
	CILIUM_CALL_IPV6_FROM_OVERLAY           = CILIUM_CALL_IPV6_FROM_LXC
	CILIUM_CALL_IPV4_TO_LXC_POLICY_ONLY     = 11
	CILIUM_CALL_IPV4_TO_HOST_POLICY_ONLY    = CILIUM_CALL_IPV4_TO_LXC_POLICY_ONLY
	CILIUM_CALL_IPV6_TO_LXC_POLICY_ONLY     = 12
	CILIUM_CALL_IPV6_TO_HOST_POLICY_ONLY    = CILIUM_CALL_IPV6_TO_LXC_POLICY_ONLY
	CILIUM_CALL_IPV4_TO_ENDPOINT            = 13
	CILIUM_CALL_IPV6_TO_ENDPOINT            = 14
	CILIUM_CALL_IPV4_NODEPORT_NAT_EGRESS    = 15
	CILIUM_CALL_IPV6_NODEPORT_NAT_EGRESS    = 16
	CILIUM_CALL_IPV4_NODEPORT_REVNAT        = 17
	CILIUM_CALL_IPV6_NODEPORT_REVNAT        = 18
	CILIUM_CALL_IPV4_NODEPORT_NAT_FWD       = 19
	CILIUM_CALL_IPV4_NODEPORT_DSR           = 20
	CILIUM_CALL_IPV6_NODEPORT_DSR           = 21
	CILIUM_CALL_IPV4_FROM_HOST              = 22
	CILIUM_CALL_IPV6_FROM_HOST              = 23
	CILIUM_CALL_IPV6_NODEPORT_NAT_FWD       = 24
	CILIUM_CALL_IPV4_FROM_LXC_CONT          = 25
	CILIUM_CALL_IPV6_FROM_LXC_CONT          = 26
	CILIUM_CALL_IPV4_CT_INGRESS             = 27
	CILIUM_CALL_IPV4_CT_INGRESS_POLICY_ONLY = 28
	CILIUM_CALL_IPV4_CT_EGRESS              = 29
	CILIUM_CALL_IPV6_CT_INGRESS             = 30
	CILIUM_CALL_IPV6_CT_INGRESS_POLICY_ONLY = 31
	CILIUM_CALL_IPV6_CT_EGRESS              = 32
	CILIUM_CALL_SRV6_ENCAP                  = 33
	CILIUM_CALL_SRV6_DECAP                  = 34
	// Unused CILIUM_CALL_SRV6_REPLY		35
	CILIUM_CALL_IPV4_NODEPORT_NAT_INGRESS = 36
	CILIUM_CALL_IPV6_NODEPORT_NAT_INGRESS = 37
	CILIUM_CALL_IPV4_NODEPORT_SNAT_FWD    = 38
	CILIUM_CALL_IPV6_NODEPORT_SNAT_FWD    = 39
	// Unused CILIUM_CALL_IPV4_NODEPORT_DSR_INGRESS	40
	// Unused CILIUM_CALL_IPV6_NODEPORT_DSR_INGRESS	41
	CILIUM_CALL_IPV4_INTER_CLUSTER_REVSNAT = 42
	CILIUM_CALL_IPV4_CONT_FROM_HOST        = 43
	CILIUM_CALL_IPV4_CONT_FROM_NETDEV      = 44
	CILIUM_CALL_IPV6_CONT_FROM_HOST        = 45
	CILIUM_CALL_IPV6_CONT_FROM_NETDEV      = 46
	CILIUM_CALL_IPV4_NO_SERVICE            = 47
	CILIUM_CALL_IPV6_NO_SERVICE            = 48
	CILIUM_CALL_MULTICAST_EP_DELIVERY      = 49
	CILIUM_CALL_SIZE                       = 50
)

// git grep -A2 "section_tail.*CILIUM_CALL_"
// Conntrack defines were manually referenced.
var TailCallMap = map[int]string{
	CILIUM_CALL_DROP_NOTIFY:                 "__send_drop_notify",
	CILIUM_CALL_ERROR_NOTIFY:                "XXX", // this one doesn't seem to be referenced.
	CILIUM_CALL_HANDLE_ICMP6_NS:             "tail_icmp6_handle_ns",
	CILIUM_CALL_SEND_ICMP6_TIME_EXCEEDED:    "tail_icmp6_send_time_exceeded",
	CILIUM_CALL_ARP:                         "tail_handle_arp",
	CILIUM_CALL_IPV4_FROM_LXC:               "tail_handle_ipv4",
	CILIUM_CALL_IPV46_RFC8215:               "tail_nat_ipv46",
	CILIUM_CALL_IPV64_RFC8215:               "tail_nat_ipv64",
	CILIUM_CALL_IPV6_FROM_LXC:               "tail_handle_ipv6",
	CILIUM_CALL_IPV4_TO_LXC_POLICY_ONLY:     "tail_ipv4_policy",
	CILIUM_CALL_IPV6_TO_LXC_POLICY_ONLY:     "tail_ipv6_policy",
	CILIUM_CALL_IPV4_TO_ENDPOINT:            "tail_ipv4_to_endpoint",
	CILIUM_CALL_IPV6_TO_ENDPOINT:            "tail_ipv6_to_endpoint",
	CILIUM_CALL_IPV4_NODEPORT_NAT_EGRESS:    "tail_nodeport_nat_egress_ipv4",
	CILIUM_CALL_IPV6_NODEPORT_NAT_EGRESS:    "tail_nodeport_nat_egress_ipv6",
	CILIUM_CALL_IPV4_NODEPORT_REVNAT:        "tail_nodeport_rev_dnat_ingress_ipv4",
	CILIUM_CALL_IPV6_NODEPORT_REVNAT:        "tail_nodeport_rev_dnat_ingress_ipv6",
	CILIUM_CALL_IPV4_NODEPORT_NAT_FWD:       "tail_handle_nat_fwd_ipv6",
	CILIUM_CALL_IPV4_NODEPORT_DSR:           "tail_nodeport_ipv4_dsr",
	CILIUM_CALL_IPV6_NODEPORT_DSR:           "tail_nodeport_ipv6_dsr",
	CILIUM_CALL_IPV4_FROM_HOST:              "tail_handle_ipv4_from_host",
	CILIUM_CALL_IPV6_FROM_HOST:              "tail_handle_ipv6_from_host",
	CILIUM_CALL_IPV6_NODEPORT_NAT_FWD:       "tail_handle_nat_fwd_ipv6",
	CILIUM_CALL_IPV4_FROM_LXC_CONT:          "tail_handle_ipv4_cont",
	CILIUM_CALL_IPV6_FROM_LXC_CONT:          "tail_handle_ipv6_cont",
	CILIUM_CALL_IPV4_CT_INGRESS:             "tail_ipv4_ct_ingress",
	CILIUM_CALL_IPV4_CT_INGRESS_POLICY_ONLY: "tail_ipv4_ct_ingress_policy_only",
	CILIUM_CALL_IPV4_CT_EGRESS:              "tail_ipv4_ct_egress",
	CILIUM_CALL_IPV6_CT_INGRESS:             "tail_ipv6_ct_ingress",
	CILIUM_CALL_IPV6_CT_INGRESS_POLICY_ONLY: "tail_ipv6_ct_ingress_policy_only",
	CILIUM_CALL_IPV6_CT_EGRESS:              "tail_ipv6_ct_egress",
	CILIUM_CALL_SRV6_ENCAP:                  "tail_srv6_encap",
	CILIUM_CALL_SRV6_DECAP:                  "tail_srv6_decap",
	CILIUM_CALL_IPV4_NODEPORT_NAT_INGRESS:   "tail_nodeport_nat_ingress_ipv4",
	CILIUM_CALL_IPV6_NODEPORT_NAT_INGRESS:   "tail_nodeport_nat_ingress_ipv6",
	CILIUM_CALL_IPV4_NODEPORT_SNAT_FWD:      "tail_handle_snat_fwd_ipv4",
	CILIUM_CALL_IPV6_NODEPORT_SNAT_FWD:      "tail_handle_snat_fwd_ipv6",
	CILIUM_CALL_IPV4_INTER_CLUSTER_REVSNAT:  "tail_handle_inter_cluster_revsnat",
	CILIUM_CALL_IPV4_CONT_FROM_HOST:         "tail_handle_ipv4_cont_from_host",
	CILIUM_CALL_IPV4_CONT_FROM_NETDEV:       "tail_handle_ipv4_from_netdev",
	CILIUM_CALL_IPV6_CONT_FROM_HOST:         "tail_handle_ipv6_cont_from_host",
	CILIUM_CALL_IPV6_CONT_FROM_NETDEV:       "tail_handle_ipv6_from_netdev",
	CILIUM_CALL_IPV4_NO_SERVICE:             "tail_no_service_ipv4",
	CILIUM_CALL_IPV6_NO_SERVICE:             "tail_no_service_ipv6",
	CILIUM_CALL_MULTICAST_EP_DELIVERY:       "tail_mcast_ep_delivery",
}
