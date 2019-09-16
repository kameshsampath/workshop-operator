// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageInfo) DeepCopyInto(out *PackageInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageInfo.
func (in *PackageInfo) DeepCopy() *PackageInfo {
	if in == nil {
		return nil
	}
	out := new(PackageInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Workshop) DeepCopyInto(out *Workshop) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Workshop.
func (in *Workshop) DeepCopy() *Workshop {
	if in == nil {
		return nil
	}
	out := new(Workshop)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Workshop) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopList) DeepCopyInto(out *WorkshopList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Workshop, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopList.
func (in *WorkshopList) DeepCopy() *WorkshopList {
	if in == nil {
		return nil
	}
	out := new(WorkshopList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WorkshopList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopProject) DeepCopyInto(out *WorkshopProject) {
	*out = *in
	if in.Prefixes != nil {
		in, out := &in.Prefixes, &out.Prefixes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopProject.
func (in *WorkshopProject) DeepCopy() *WorkshopProject {
	if in == nil {
		return nil
	}
	out := new(WorkshopProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopSpec) DeepCopyInto(out *WorkshopSpec) {
	*out = *in
	in.Project.DeepCopyInto(&out.Project)
	out.User = in.User
	in.Stack.DeepCopyInto(&out.Stack)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopSpec.
func (in *WorkshopSpec) DeepCopy() *WorkshopSpec {
	if in == nil {
		return nil
	}
	out := new(WorkshopSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopStack) DeepCopyInto(out *WorkshopStack) {
	*out = *in
	if in.Community != nil {
		in, out := &in.Community, &out.Community
		*out = make([]PackageInfo, len(*in))
		copy(*out, *in)
	}
	if in.RedHat != nil {
		in, out := &in.RedHat, &out.RedHat
		*out = make([]PackageInfo, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopStack.
func (in *WorkshopStack) DeepCopy() *WorkshopStack {
	if in == nil {
		return nil
	}
	out := new(WorkshopStack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopStatus) DeepCopyInto(out *WorkshopStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopStatus.
func (in *WorkshopStatus) DeepCopy() *WorkshopStatus {
	if in == nil {
		return nil
	}
	out := new(WorkshopStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkshopUser) DeepCopyInto(out *WorkshopUser) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkshopUser.
func (in *WorkshopUser) DeepCopy() *WorkshopUser {
	if in == nil {
		return nil
	}
	out := new(WorkshopUser)
	in.DeepCopyInto(out)
	return out
}
