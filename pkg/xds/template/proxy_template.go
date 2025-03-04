package template

import (
	kuma_mesh "github.com/kumahq/kuma/api/mesh/v1alpha1"
	core_mesh "github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/xds/generator"
)

var (
	DefaultProxyTemplate = &kuma_mesh.ProxyTemplate{
		Conf: &kuma_mesh.ProxyTemplate_Conf{
			Imports: []string{
				core_mesh.ProfileDefaultProxy,
			},
		},
	}

	IngressProxyTemplate = &kuma_mesh.ProxyTemplate{
		Conf: &kuma_mesh.ProxyTemplate_Conf{
			Imports: []string{
				generator.IngressProxy,
			},
		},
	}
)
