package upcloud

import (
	"testing"

	"github.com/UpCloudLtd/terraform-provider-upcloud/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUpcloudKubernetes(t *testing.T) {
	testDataS1 := utils.ReadTestDataFile(t, "testdata/upcloud_kubernetes/kubernetes_s1.tf")
	testDataS2 := utils.ReadTestDataFile(t, "testdata/upcloud_kubernetes/kubernetes_s2.tf")

	cName := "upcloud_kubernetes_cluster.main"
	g1Name := "upcloud_kubernetes_node_group.g1"
	g2Name := "upcloud_kubernetes_node_group.g2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataS1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(cName, "control_plane_ip_filter.*", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(cName, "name", "tf-acc-test-uks"),
					resource.TestCheckResourceAttr(cName, "version", "1.27"),
					resource.TestCheckResourceAttr(cName, "zone", "fi-hel2"),
					resource.TestCheckResourceAttr(g1Name, "name", "small"),
					resource.TestCheckResourceAttr(g2Name, "name", "medium"),
					resource.TestCheckResourceAttr(g1Name, "anti_affinity", "true"),
					resource.TestCheckResourceAttr(g2Name, "anti_affinity", "false"),
					resource.TestCheckResourceAttr(g1Name, "node_count", "2"),
					resource.TestCheckResourceAttr(g2Name, "node_count", "1"),
					resource.TestCheckResourceAttr(g1Name, "ssh_keys.#", "1"),
					resource.TestCheckResourceAttr(g2Name, "ssh_keys.#", "1"),
					resource.TestCheckResourceAttr(g1Name, "labels.%", "2"),
					resource.TestCheckResourceAttr(g2Name, "labels.%", "2"),
					resource.TestCheckResourceAttr(g1Name, "labels.env", "dev"),
					resource.TestCheckResourceAttr(g1Name, "labels.managedBy", "tf"),
					resource.TestCheckTypeSetElemAttr(g1Name, "ssh_keys.*", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testt_test"),
					resource.TestCheckTypeSetElemNestedAttrs(g1Name, "kubelet_args.*", map[string]string{
						"key":   "log-flush-frequency",
						"value": "5s",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(g1Name, "taint.*", map[string]string{
						"effect": "NoExecute",
						"key":    "taintKey",
						"value":  "taintValue",
					}),
					resource.TestCheckResourceAttr(g1Name, "utility_network_access", "true"),
					resource.TestCheckResourceAttr(g2Name, "utility_network_access", "false"),
				),
			},
			{
				Config:            testDataS2,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cName, "control_plane_ip_filter.#", "0"),
					resource.TestCheckResourceAttr(cName, "version", "1.27"),
					resource.TestCheckResourceAttr(g1Name, "node_count", "1"),
					resource.TestCheckResourceAttr(g2Name, "node_count", "2"),
				),
			},
		},
	})
}

func TestAccUpcloudKubernetes_labels(t *testing.T) {
	testDataS1 := utils.ReadTestDataFile(t, "testdata/upcloud_kubernetes/kubernetes_labels_s1.tf")
	testDataS2 := utils.ReadTestDataFile(t, "testdata/upcloud_kubernetes/kubernetes_labels_s2.tf")
	testDataS3 := utils.ReadTestDataFile(t, "testdata/upcloud_kubernetes/kubernetes_labels_s3.tf")

	cluster := "upcloud_kubernetes_cluster.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataS1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(cluster, "control_plane_ip_filter.*", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(cluster, "name", "tf-acc-test-k8s-labels-cluster"),
					resource.TestCheckResourceAttr(cluster, "zone", "de-fra1"),
					resource.TestCheckResourceAttr(cluster, "labels.%", "1"),
					resource.TestCheckResourceAttr(cluster, "labels.test", "terraform-provider-acceptance-test"),
				),
			},
			{
				Config: testDataS2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cluster, "name", "tf-acc-test-k8s-labels-cluster"),
					resource.TestCheckResourceAttr(cluster, "labels.%", "2"),
					resource.TestCheckResourceAttr(cluster, "labels.test", "terraform-provider-acceptance-test"),
					resource.TestCheckResourceAttr(cluster, "labels.managed-by", "terraform"),
				),
			},
			{
				Config: testDataS3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cluster, "name", "tf-acc-test-k8s-labels-cluster"),
					resource.TestCheckResourceAttr(cluster, "labels.%", "0"),
				),
			},
		},
	})
}
