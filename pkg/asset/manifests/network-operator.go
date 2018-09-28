package manifests

import (
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"

	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
)

// networkOperator generates the network-operator-*.yml files
type networkOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
}

var _ asset.Asset = (*networkOperator)(nil)

// Name returns a human friendly name for the operator
func (no *networkOperator) Name() string {
	return "Network Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// networkOperator asset.
func (no *networkOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		no.installConfigAsset,
	}
}

// Generate generates the network-operator-config.yml and network-operator-manifest.yml files
func (no *networkOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(no.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	no.installConfig = ic

	// installconfig is ready, we can create the core config from it now
	netConfig, err := no.netConfig()
	if err != nil {
		return nil, err
	}

	netManifest, err := no.manifest()
	if err != nil {
		return nil, err
	}
	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: "network-operator-config.yml",
				Data: netConfig,
			},
			{
				Name: "network-operator-manifests.yml",
				Data: netManifest,
			},
		},
	}
	return state, nil
}

func (no *networkOperator) netConfig() ([]byte, error) {
	networkConfig := netopv1.NetworkConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: netopv1.SchemeGroupVersion.String(),
			Kind:       "NetworkConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: netopv1.NetworkConfigSpec{
			ServiceNetwork: no.installConfig.Networking.ServiceCIDR.String(),
			ClusterNetworks: []netopv1.ClusterNetwork{
				{
					CIDR:             no.installConfig.Networking.PodCIDR.String(),
					HostSubnetLength: 9,
				},
			},
			DefaultNetwork: netopv1.DefaultNetworkDefinition{
				Type: netopv1.NetworkTypeOpenshiftSDN,
				OpenshiftSDNConfig: &netopv1.OpenshiftSDNConfig{
					Mode: netopv1.SDNModePolicy,
				},
			},
		},
	}

	return yaml.Marshal(networkConfig)
}

func (no *networkOperator) manifest() ([]byte, error) {
	return []byte(""), nil
}
