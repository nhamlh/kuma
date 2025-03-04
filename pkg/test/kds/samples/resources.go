package samples

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	mesh_proto "github.com/kumahq/kuma/api/mesh/v1alpha1"
	system_proto "github.com/kumahq/kuma/api/system/v1alpha1"
)

var (
	Mesh1 = &mesh_proto.Mesh{
		Mtls: &mesh_proto.Mesh_Mtls{
			EnabledBackend: "ca-1",
			Backends: []*mesh_proto.CertificateAuthorityBackend{
				{
					Name: "ca-1",
					Type: "builtin",
				},
			},
		},
	}
	Mesh2 = &mesh_proto.Mesh{
		Mtls: &mesh_proto.Mesh_Mtls{
			EnabledBackend: "ca-2",
			Backends: []*mesh_proto.CertificateAuthorityBackend{
				{
					Name: "ca-2",
					Type: "builtin",
				},
			},
		},
	}
	FaultInjection = &mesh_proto.FaultInjection{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
				"tag0":    "version0",
				"tag1":    "version1",
				"tag2":    "version2",
				"tag3":    "version3",
				"tag4":    "version4",
				"tag5":    "version5",
				"tag6":    "version6",
				"tag7":    "version7",
				"tag8":    "version8",
				"tag9":    "version9",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.FaultInjection_Conf{
			Abort: &mesh_proto.FaultInjection_Conf_Abort{
				Percentage: &wrapperspb.DoubleValue{Value: 90},
				HttpStatus: &wrapperspb.UInt32Value{Value: 404},
			},
		},
	}
	Dataplane = &mesh_proto.Dataplane{
		Networking: &mesh_proto.Dataplane_Networking{
			Address: "192.168.0.1",
			Inbound: []*mesh_proto.Dataplane_Networking_Inbound{{
				Port: 1212,
				Tags: map[string]string{
					mesh_proto.ZoneTag:    "kuma-1",
					mesh_proto.ServiceTag: "backend",
				},
			}},
			Outbound: []*mesh_proto.Dataplane_Networking_Outbound{
				{
					Port: 1213,
					Tags: map[string]string{
						mesh_proto.ServiceTag:  "web",
						mesh_proto.ProtocolTag: "http",
					},
				},
			},
		},
	}
	DataplaneInsight = &mesh_proto.DataplaneInsight{
		MTLS: &mesh_proto.DataplaneInsight_MTLS{
			CertificateRegenerations: 3,
		},
	}
	Ingress = &mesh_proto.Dataplane{
		Networking: &mesh_proto.Dataplane_Networking{
			Ingress: &mesh_proto.Dataplane_Networking_Ingress{
				AvailableServices: []*mesh_proto.Dataplane_Networking_Ingress_AvailableService{{
					Tags: map[string]string{
						"service": "backend",
					}},
				},
			},
			Address: "192.168.0.1",
		},
	}
	ZoneIngress = &mesh_proto.ZoneIngress{
		Networking: &mesh_proto.ZoneIngress_Networking{
			Address:           "127.0.0.1",
			Port:              80,
			AdvertisedAddress: "192.168.0.1",
			AdvertisedPort:    10001,
		},
		AvailableServices: []*mesh_proto.ZoneIngress_AvailableService{{
			Tags: map[string]string{
				"service": "backend",
			}},
		},
	}
	ZoneIngressInsight = &mesh_proto.ZoneIngressInsight{
		Subscriptions: []*mesh_proto.DiscoverySubscription{{
			Id: "1",
		}},
	}
	ExternalService = &mesh_proto.ExternalService{
		Networking: &mesh_proto.ExternalService_Networking{
			Address: "192.168.0.1",
		},
		Tags: map[string]string{
			mesh_proto.ZoneTag:    "kuma-1",
			mesh_proto.ServiceTag: "backend",
		},
	}
	CircuitBreaker = &mesh_proto.CircuitBreaker{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.CircuitBreaker_Conf{
			Detectors: &mesh_proto.CircuitBreaker_Conf_Detectors{},
		},
	}
	HealthCheck = &mesh_proto.HealthCheck{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.HealthCheck_Conf{
			Interval: &durationpb.Duration{Seconds: 5},
			Timeout:  &durationpb.Duration{Seconds: 7},
		},
	}
	TrafficLog = &mesh_proto.TrafficLog{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.TrafficLog_Conf{
			Backend: "logging-backend",
		},
	}
	TrafficPermission = &mesh_proto.TrafficPermission{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"kuma.io/service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"kuma.io/service": "*",
			},
		}},
	}
	TrafficRoute = &mesh_proto.TrafficRoute{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.TrafficRoute_Conf{
			Split: []*mesh_proto.TrafficRoute_Split{{
				Weight: &wrapperspb.UInt32Value{
					Value: 10,
				},
				Destination: map[string]string{
					"version": "v2",
				},
			}},
		},
	}
	TrafficTrace = &mesh_proto.TrafficTrace{
		Selectors: []*mesh_proto.Selector{{
			Match: map[string]string{"serivce": "*"},
		}},
		Conf: &mesh_proto.TrafficTrace_Conf{
			Backend: "tracing-backend",
		},
	}
	ProxyTemplate = &mesh_proto.ProxyTemplate{
		Selectors: []*mesh_proto.Selector{{
			Match: map[string]string{"serivce": "*"},
		}},
		Conf: &mesh_proto.ProxyTemplate_Conf{
			Imports: []string{"default-kuma-profile"},
		},
	}
	Retry = &mesh_proto.Retry{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.Retry_Conf{
			Http: &mesh_proto.Retry_Conf_Http{
				NumRetries: &wrapperspb.UInt32Value{
					Value: 5,
				},
				PerTryTimeout: &durationpb.Duration{
					Seconds: 200000000,
				},
				BackOff: &mesh_proto.Retry_Conf_BackOff{
					BaseInterval: &durationpb.Duration{
						Nanos: 200000000,
					},
					MaxInterval: &durationpb.Duration{
						Seconds: 1,
					},
				},
				RetriableStatusCodes: []uint32{500, 502},
			},
		},
	}
	Timeout = &mesh_proto.Timeout{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"service": "*",
			},
		}},
		Conf: &mesh_proto.Timeout_Conf{

			ConnectTimeout: &durationpb.Duration{
				Seconds: 5,
			},
			Tcp: &mesh_proto.Timeout_Conf_Tcp{
				IdleTimeout: &durationpb.Duration{
					Seconds: 5,
				},
			},
			Http: &mesh_proto.Timeout_Conf_Http{
				RequestTimeout: &durationpb.Duration{
					Seconds: 5,
				},
				IdleTimeout: &durationpb.Duration{
					Seconds: 5,
				},
			},
			Grpc: &mesh_proto.Timeout_Conf_Grpc{
				StreamIdleTimeout: &durationpb.Duration{
					Seconds: 5,
				},
				MaxStreamDuration: &durationpb.Duration{
					Seconds: 5,
				},
			},
		},
	}
	Secret = &system_proto.Secret{
		Data: &wrapperspb.BytesValue{Value: []byte("secret key")},
	}
	GlobalSecret = &system_proto.Secret{
		Data: &wrapperspb.BytesValue{Value: []byte("global secret key")},
	}
	Config = &system_proto.Config{
		Config: "sample config",
	}
	RateLimit = &mesh_proto.RateLimit{
		Sources: []*mesh_proto.Selector{{
			Match: map[string]string{
				"kuma.io/service": "*",
			},
		}},
		Destinations: []*mesh_proto.Selector{{
			Match: map[string]string{
				"kuma.io/service": "*",
			},
		}},
	}
)
