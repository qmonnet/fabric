// Copyright 2023 Hedgehog
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.githedgehog.com/fabric/api/meta"
	wiringapi "go.githedgehog.com/fabric/api/wiring/v1beta1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetAPI2NOSPortsFor(t *testing.T) {
	for _, tt := range []struct {
		name string
		sp   *wiringapi.SwitchProfileSpec
		sw   *wiringapi.SwitchSpec
		want map[string]string
		err  bool
	}{
		{
			name: "simple",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"M1":   {NOSName: "Management0", Management: true, OniePortName: "eth0"},
					"E1/1": {NOSName: "Ethernet0", Label: "1", Group: "1"},
					"E1/2": {NOSName: "Ethernet4", Label: "2", Group: "1"},
					"E1/3": {NOSName: "Ethernet8", Label: "3", Group: "2"},
					"E1/4": {NOSName: "Ethernet12", Label: "4", Group: "2"},
					"E1/5": {NOSName: "Ethernet16", Label: "5", Profile: "SFP28-25G"},
					"E1/6": {NOSName: "Ethernet17", Label: "6", Profile: "SFP28-25G"},
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
					"E1/9": {NOSName: "1/9", Label: "9", Profile: "QSFP28-100G", BaseNOSName: "Ethernet28"},
				},
				PortGroups: map[string]wiringapi.SwitchProfilePortGroup{
					"1": {
						NOSName: "1",
						Profile: "SFP28-25G",
					},
					"2": {
						NOSName: "2",
						Profile: "SFP28-25G",
					},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"SFP28-25G": {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "25G",
							Supported: []string{"10G", "25G"},
						},
					},
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			sw: &wiringapi.SwitchSpec{
				PortBreakouts: map[string]string{
					"E1/8": "4x25G",
					"E1/9": "2x50G",
				},
			},
			want: map[string]string{
				"E1/1":   "Ethernet0",
				"E1/2":   "Ethernet4",
				"E1/3":   "Ethernet8",
				"E1/4":   "Ethernet12",
				"E1/5":   "Ethernet16",
				"E1/6":   "Ethernet17",
				"E1/7":   "Ethernet20",
				"E1/7/1": "Ethernet20",
				"E1/8":   "Ethernet24",
				"E1/8/1": "Ethernet24",
				"E1/8/2": "Ethernet25",
				"E1/8/3": "Ethernet26",
				"E1/8/4": "Ethernet27",
				"E1/9":   "Ethernet28",
				"E1/9/1": "Ethernet28",
				"E1/9/2": "Ethernet30",
				"M1":     "Management0",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sp != nil {
				tt.sp.NOSType = meta.NOSTypeSONiCBCMVS
				tt.sp.Platform = "vs"
			}
			_, err := (&wiringapi.SwitchProfile{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test",
					Namespace: kmetav1.NamespaceDefault,
				},
				Spec: *tt.sp,
			}).Validate(t.Context(), nil, nil)
			require.NoError(t, err)

			got, err := tt.sp.GetAPI2NOSPortsFor(tt.sw)

			if tt.err {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetNOS2APIPortsFor(t *testing.T) {
	for _, tt := range []struct {
		name string
		sp   *wiringapi.SwitchProfileSpec
		sw   *wiringapi.SwitchSpec
		want map[string]string
		err  bool
	}{
		{
			name: "simple",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"M1":   {NOSName: "Management0", Management: true, OniePortName: "eth0"},
					"E1/1": {NOSName: "Ethernet0", Label: "1", Group: "1"},
					"E1/2": {NOSName: "Ethernet4", Label: "2", Group: "1"},
					"E1/3": {NOSName: "Ethernet8", Label: "3", Group: "2"},
					"E1/4": {NOSName: "Ethernet12", Label: "4", Group: "2"},
					"E1/5": {NOSName: "Ethernet16", Label: "5", Profile: "SFP28-25G"},
					"E1/6": {NOSName: "Ethernet17", Label: "6", Profile: "SFP28-25G"},
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
					"E1/9": {NOSName: "1/9", Label: "9", Profile: "QSFP28-100G", BaseNOSName: "Ethernet28"},
				},
				PortGroups: map[string]wiringapi.SwitchProfilePortGroup{
					"1": {
						NOSName: "1",
						Profile: "SFP28-25G",
					},
					"2": {
						NOSName: "2",
						Profile: "SFP28-25G",
					},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"SFP28-25G": {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "25G",
							Supported: []string{"10G", "25G"},
						},
					},
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			sw: &wiringapi.SwitchSpec{
				PortBreakouts: map[string]string{
					"E1/8": "4x25G",
					"E1/9": "2x50G",
				},
			},
			want: map[string]string{
				"Ethernet0":   "E1/1",
				"Ethernet4":   "E1/2",
				"Ethernet8":   "E1/3",
				"Ethernet12":  "E1/4",
				"Ethernet16":  "E1/5",
				"Ethernet17":  "E1/6",
				"Ethernet20":  "E1/7/1",
				"Ethernet24":  "E1/8/1",
				"Ethernet25":  "E1/8/2",
				"Ethernet26":  "E1/8/3",
				"Ethernet27":  "E1/8/4",
				"Ethernet28":  "E1/9/1",
				"Ethernet30":  "E1/9/2",
				"Management0": "M1",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sp != nil {
				tt.sp.NOSType = meta.NOSTypeSONiCBCMVS
				tt.sp.Platform = "vs"
			}
			_, err := (&wiringapi.SwitchProfile{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test",
					Namespace: kmetav1.NamespaceDefault,
				},
				Spec: *tt.sp,
			}).Validate(t.Context(), nil, nil)
			require.NoError(t, err)

			got, err := tt.sp.GetNOS2APIPortsFor(tt.sw)

			if tt.err {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetAllBreakoutNOSNames(t *testing.T) {
	for _, tt := range []struct {
		name string
		sp   *wiringapi.SwitchProfileSpec
		want map[string]bool
		err  bool
	}{
		{
			name: "simple",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			want: map[string]bool{
				"Ethernet20": true,
				"Ethernet21": true,
				"Ethernet22": true,
				"Ethernet23": true,
				"Ethernet24": true,
				"Ethernet25": true,
				"Ethernet26": true,
				"Ethernet27": true,
			},
		},
		{
			name: "simple2",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			want: map[string]bool{
				"Ethernet20": true,
				"Ethernet22": true,
				"Ethernet24": true,
				"Ethernet26": true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sp != nil {
				tt.sp.NOSType = meta.NOSTypeSONiCBCMVS
				tt.sp.Platform = "vs"
			}
			_, err := (&wiringapi.SwitchProfile{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test",
					Namespace: kmetav1.NamespaceDefault,
				},
				Spec: *tt.sp,
			}).Validate(t.Context(), nil, nil)
			require.NoError(t, err)

			got, err := tt.sp.GetAllBreakoutNOSNames()

			if tt.err {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSwitchProfileDefault(t *testing.T) {
	for _, tt := range []struct {
		name string
		sp   *wiringapi.SwitchProfileSpec
		want map[string]string
	}{
		{
			name: "simple",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"M1":    {NOSName: "Management0", Management: true, OniePortName: "eth0"},
					"E1/1":  {NOSName: "Ethernet0", Label: "1", Group: "1"},
					"E1/2":  {NOSName: "Ethernet4", Label: "2", Group: "1"},
					"E1/3":  {NOSName: "Ethernet8", Label: "3", Group: "2"},
					"E1/4":  {NOSName: "Ethernet12", Label: "4", Group: "2"},
					"E1/5":  {NOSName: "Ethernet16", Label: "5", Profile: "SFP28-25G"},
					"E1/6":  {NOSName: "Ethernet17", Label: "6", Profile: "SFP28-25G"},
					"E1/7":  {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8":  {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
					"E1/9":  {NOSName: "1/9", Label: "9", Profile: "QSFP28-100G", BaseNOSName: "Ethernet28"},
					"E1/10": {NOSName: "Ethernet42", Label: "10", Profile: "QSFP28-100G" + wiringapi.NonBreakoutPortExceptionSuffix},
				},
				PortGroups: map[string]wiringapi.SwitchProfilePortGroup{
					"1": {
						NOSName: "1",
						Profile: "SFP28-25G",
					},
					"2": {
						NOSName: "2",
						Profile: "SFP28-25G",
					},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"SFP28-25G": {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "25G",
							Supported: []string{"10G", "25G"},
						},
					},
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
					"QSFP28-100G" + wiringapi.NonBreakoutPortExceptionSuffix: {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "100G",
							Supported: []string{"10G", "25G", "50G", "100G"},
						},
					},
				},
			},
			want: map[string]string{
				wiringapi.AnnotationPorts: "6xSFP28-25G, 4xQSFP28-100G",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sp != nil {
				tt.sp.NOSType = meta.NOSTypeSONiCBCMVS
				tt.sp.Platform = "vs"
			}
			_, err := (&wiringapi.SwitchProfile{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test",
					Namespace: kmetav1.NamespaceDefault,
				},
				Spec: *tt.sp,
			}).Validate(t.Context(), nil, nil)
			require.NoError(t, err)

			sp := &wiringapi.SwitchProfile{
				Spec: *tt.sp,
			}
			sp.Default()

			require.Equal(t, tt.want, sp.Annotations)
		})
	}
}

func TestComparePortNames(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{
			name: "simple-E",
			a:    "E1/1",
			b:    "E1/2",
			want: -1,
		},
		{
			name: "simple-E-same",
			a:    "E1/1",
			b:    "E1/1",
			want: 0,
		},
		{
			name: "simple-M",
			a:    "M1",
			b:    "M2",
			want: -1,
		},
		{
			name: "simple-M-same",
			a:    "M1",
			b:    "M1",
			want: 0,
		},
		{
			name: "simple-ME",
			a:    "E1/1",
			b:    "M1",
			want: 1,
		},
		{
			name: "simple-EB",
			a:    "E1/1",
			b:    "E1/1/1",
			want: -1,
		},
		{
			name: "simple-BB",
			a:    "E1/1/1",
			b:    "E1/1/2",
			want: -1,
		},
		{
			name: "simple-MB",
			a:    "M1",
			b:    "E1/1/1",
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wiringapi.ComparePortNames(tt.a, tt.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetAvailableAPIPorts(t *testing.T) {
	for _, tt := range []struct {
		name string
		sp   *wiringapi.SwitchProfileSpec
		sw   *wiringapi.SwitchSpec
		want map[string]bool
		err  bool
	}{
		{
			name: "simple",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"M1":   {NOSName: "Management0", Management: true, OniePortName: "eth0"},
					"E1/1": {NOSName: "Ethernet0", Label: "1", Group: "1"},
					"E1/2": {NOSName: "Ethernet4", Label: "2", Group: "1"},
					"E1/3": {NOSName: "Ethernet8", Label: "3", Group: "2"},
					"E1/4": {NOSName: "Ethernet12", Label: "4", Group: "2"},
					"E1/5": {NOSName: "Ethernet16", Label: "5", Profile: "SFP28-25G"},
					"E1/6": {NOSName: "Ethernet17", Label: "6", Profile: "SFP28-25G"},
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
					"E1/9": {NOSName: "1/9", Label: "9", Profile: "QSFP28-100G", BaseNOSName: "Ethernet28"},
				},
				PortGroups: map[string]wiringapi.SwitchProfilePortGroup{
					"1": {
						NOSName: "1",
						Profile: "SFP28-25G",
					},
					"2": {
						NOSName: "2",
						Profile: "SFP28-25G",
					},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"SFP28-25G": {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "25G",
							Supported: []string{"10G", "25G"},
						},
					},
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			sw: &wiringapi.SwitchSpec{
				PortBreakouts: map[string]string{
					"E1/8": "4x25G",
					"E1/9": "2x50G",
				},
			},
			want: map[string]bool{
				"E1/1":   true,
				"E1/2":   true,
				"E1/3":   true,
				"E1/4":   true,
				"E1/5":   true,
				"E1/6":   true,
				"E1/7/1": true,
				"E1/8/1": true,
				"E1/8/2": true,
				"E1/8/3": true,
				"E1/8/4": true,
				"E1/9/1": true,
				"E1/9/2": true,
			},
		},
		{
			name: "simple-no-switch",
			sp: &wiringapi.SwitchProfileSpec{
				DisplayName:   "Test",
				SwitchSilicon: "test",
				Ports: map[string]wiringapi.SwitchProfilePort{
					"M1":   {NOSName: "Management0", Management: true, OniePortName: "eth0"},
					"E1/1": {NOSName: "Ethernet0", Label: "1", Group: "1"},
					"E1/2": {NOSName: "Ethernet4", Label: "2", Group: "1"},
					"E1/3": {NOSName: "Ethernet8", Label: "3", Group: "2"},
					"E1/4": {NOSName: "Ethernet12", Label: "4", Group: "2"},
					"E1/5": {NOSName: "Ethernet16", Label: "5", Profile: "SFP28-25G"},
					"E1/6": {NOSName: "Ethernet17", Label: "6", Profile: "SFP28-25G"},
					"E1/7": {NOSName: "1/7", Label: "7", Profile: "QSFP28-100G", BaseNOSName: "Ethernet20"},
					"E1/8": {NOSName: "1/8", Label: "8", Profile: "QSFP28-100G", BaseNOSName: "Ethernet24"},
					"E1/9": {NOSName: "1/9", Label: "9", Profile: "QSFP28-100G", BaseNOSName: "Ethernet28"},
				},
				PortGroups: map[string]wiringapi.SwitchProfilePortGroup{
					"1": {
						NOSName: "1",
						Profile: "SFP28-25G",
					},
					"2": {
						NOSName: "2",
						Profile: "SFP28-25G",
					},
				},
				PortProfiles: map[string]wiringapi.SwitchProfilePortProfile{
					"SFP28-25G": {
						Speed: &wiringapi.SwitchProfilePortProfileSpeed{
							Default:   "25G",
							Supported: []string{"10G", "25G"},
						},
					},
					"QSFP28-100G": {
						Breakout: &wiringapi.SwitchProfilePortProfileBreakout{
							Default: "1x100G",
							Supported: map[string]wiringapi.SwitchProfilePortProfileBreakoutMode{
								"1x100G": {Offsets: []string{"0"}},
								"1x40G":  {Offsets: []string{"0"}},
								"2x50G":  {Offsets: []string{"0", "2"}},
								"1x50G":  {Offsets: []string{"0"}},
								"4x25G":  {Offsets: []string{"0", "1", "2", "3"}},
								"4x10G":  {Offsets: []string{"0", "1", "2", "3"}},
								"1x25G":  {Offsets: []string{"0"}},
								"1x10G":  {Offsets: []string{"0"}},
							},
						},
					},
				},
			},
			sw: nil,
			want: map[string]bool{
				"E1/1":   true,
				"E1/2":   true,
				"E1/3":   true,
				"E1/4":   true,
				"E1/5":   true,
				"E1/6":   true,
				"E1/7/1": true,
				"E1/8/1": true,
				"E1/9/1": true,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sp != nil {
				tt.sp.NOSType = meta.NOSTypeSONiCBCMVS
				tt.sp.Platform = "vs"
			}

			_, err := (&wiringapi.SwitchProfile{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:      "test",
					Namespace: kmetav1.NamespaceDefault,
				},
				Spec: *tt.sp,
			}).Validate(t.Context(), nil, nil)
			require.NoError(t, err)

			got, err := tt.sp.GetAvailableAPIPorts(tt.sw)

			if tt.err {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
