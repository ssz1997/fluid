/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubeclient

import (
	"testing"

	"github.com/fluid-cloudnative/fluid/pkg/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// Use fake client because of it will be maintained in the long term
// due to https://github.com/kubernetes-sigs/controller-runtime/pull/1101
func TestIsPersistentVolumeExist(t *testing.T) {

	testPVInputs := []*v1.PersistentVolume{&v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{Name: "notCreatedByFluid"},
		Spec:       v1.PersistentVolumeSpec{},
	}, &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{Name: "createdByFluid", Annotations: common.ExpectedFluidAnnotations},
		Spec:       v1.PersistentVolumeSpec{},
	}}

	testPVs := []runtime.Object{}

	for _, pv := range testPVInputs {
		testPVs = append(testPVs, pv.DeepCopy())
	}

	client := fake.NewFakeClientWithScheme(testScheme, testPVs...)

	type args struct {
		name        string
		annotations map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "volume doesn't exist",
			args: args{
				name:        "notExist",
				annotations: map[string]string{},
			},
			want: false,
		},
		{
			name: "volume is not created by fluid",
			args: args{
				name:        "notCreatedByFluid",
				annotations: map[string]string{},
			},
			want: false,
		},
		{
			name: "volume is created by fluid",
			args: args{
				name:        "createdByFluid",
				annotations: common.ExpectedFluidAnnotations,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := IsPersistentVolumeExist(client, tt.args.name, tt.args.annotations); got != tt.want {
				t.Errorf("testcase %v IsPersistentVolumeExist() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}

}