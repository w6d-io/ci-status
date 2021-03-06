/*
Copyright 2020 WILDCARD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Created on 28/01/2021
*/
package tekton_test

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis/duck/v1beta1"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tkn "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/w6d-io/ci-status/internal/tekton"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Tekton", func() {
	Context("GetWatch", func() {
		When("resource does not exist", func() {
			It("returns nil for pods", func() {
				t := &tekton.Tekton{
					PipelineRun: tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
				}
				Expect(t.GetWatch("pods", "test")).To(BeNil())
			})
		})
	})
	Context("GetTask", func() {
		When("resource does not exist", func() {
			It("returns nil for pods", func() {
				t := &tekton.Tekton{
					PipelineRun: tekton.PipelineRunPayload{
						NamespacedName: types.NamespacedName{
							Name:      "pod-test-1-1",
							Namespace: "default",
						},
					},
				}
				pr := &tkn.PipelineRun{
					Status: tkn.PipelineRunStatus{
						PipelineRunStatusFields: tkn.PipelineRunStatusFields{
							StartTime: &metav1.Time{
								Time: time.Now().AddDate(0, 0, -1),
							},
							CompletionTime: &metav1.Time{
								Time: time.Now(),
							},
						},
					},
				}
				Expect(t.PipelineRun.GetTasks(pr)).To(BeNil())
			})
		})
	})
	Context("Tasks sort", func() {
		It("Swap", func() {
			ts := tekton.Tasks{
				{
					Name: "task1",
				},
				{
					Name:      "task2",
					StartTime: 1,
				},
			}
			Expect(ts.Len()).To(Equal(2))
			Expect(ts.Less(0, 1)).To(Equal(true))
			ts.Swap(0, 1)
			ts[1].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, -1),
			}
			ts.Swap(0, 1)
			ts[0].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, 0),
			}
			ts.Swap(0, 1)
			ts[1].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, 0),
			}
			ts[0].StartTimeRaw = &metav1.Time{
				Time: time.Now().AddDate(0, 0, -1),
			}
			ts.Swap(0, 1)
		})
	})
	Context("get condition", func() {
		It("get none", func() {
			c := v1beta1.Conditions{}
			status, reason := tekton.Condition(c)
			Expect(status).To(Equal("---"))
			Expect(reason).To(Equal(""))
		})
		It("get a status Cancelled", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionFalse,
					Reason: "PipelineRunCancelled",
				},
			}
			status, reason := tekton.Condition(c)
			Expect(status).To(Equal("Cancelled"))
			Expect(reason).To(Equal("Pipeline cancelled"))
		})
		It("get a status Cancelled", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionFalse,
					Reason: "PipelineRunTimeout",
				},
			}
			status, reason := tekton.Condition(c)
			Expect(status).To(Equal("Failed"))
			Expect(reason).To(Equal("Pipeline timeout"))
		})
		It("get a status Running", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionUnknown,
				},
			}
			status, reason := tekton.Condition(c)
			Expect(status).To(Equal("Running"))
			Expect(reason).To(Equal(""))
		})
		It("get a status Succeeded", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionTrue,
				},
			}
			status, reason := tekton.Condition(c)
			Expect(status).To(Equal("Succeeded"))
			Expect(reason).To(Equal(""))
		})
	})
	Context("check condition terminated", func() {
		It("get false because empty", func() {
			c := v1beta1.Conditions{}
			Expect(tekton.IsTerminated(c)).To(Equal(false))
		})
		It("get true because condition is false", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionFalse,
					Reason: "PipelineRunCancelled",
				},
			}
			Expect(tekton.IsTerminated(c)).To(Equal(true))
		})
		It("get false because condition is unknown", func() {
			c := v1beta1.Conditions{
				{
					Status: corev1.ConditionUnknown,
					Reason: "PipelineRunCancelled",
				},
			}
			Expect(tekton.IsTerminated(c)).To(Equal(false))
		})
	})
})
