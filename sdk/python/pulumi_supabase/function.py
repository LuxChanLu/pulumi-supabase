# coding=utf-8
# *** WARNING: this file was generated by Pulumi SDK Generator. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities
from ._enums import *

__all__ = ['FunctionArgs', 'Function']

@pulumi.input_type
class FunctionArgs:
    def __init__(__self__, *,
                 body: pulumi.Input[str],
                 name: pulumi.Input[str],
                 project_id: pulumi.Input[str],
                 slug: pulumi.Input[str],
                 verify_jwt: Optional[pulumi.Input[bool]] = None):
        """
        The set of arguments for constructing a Function resource.
        :param pulumi.Input[str] body: Body of the functino
        :param pulumi.Input[str] name: Name of the function
        :param pulumi.Input[str] project_id: ID of the project
        :param pulumi.Input[str] slug: Slug of the function
        :param pulumi.Input[bool] verify_jwt: Verify JWT before running
        """
        pulumi.set(__self__, "body", body)
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "project_id", project_id)
        pulumi.set(__self__, "slug", slug)
        if verify_jwt is None:
            verify_jwt = False
        if verify_jwt is not None:
            pulumi.set(__self__, "verify_jwt", verify_jwt)

    @property
    @pulumi.getter
    def body(self) -> pulumi.Input[str]:
        """
        Body of the functino
        """
        return pulumi.get(self, "body")

    @body.setter
    def body(self, value: pulumi.Input[str]):
        pulumi.set(self, "body", value)

    @property
    @pulumi.getter
    def name(self) -> pulumi.Input[str]:
        """
        Name of the function
        """
        return pulumi.get(self, "name")

    @name.setter
    def name(self, value: pulumi.Input[str]):
        pulumi.set(self, "name", value)

    @property
    @pulumi.getter(name="projectId")
    def project_id(self) -> pulumi.Input[str]:
        """
        ID of the project
        """
        return pulumi.get(self, "project_id")

    @project_id.setter
    def project_id(self, value: pulumi.Input[str]):
        pulumi.set(self, "project_id", value)

    @property
    @pulumi.getter
    def slug(self) -> pulumi.Input[str]:
        """
        Slug of the function
        """
        return pulumi.get(self, "slug")

    @slug.setter
    def slug(self, value: pulumi.Input[str]):
        pulumi.set(self, "slug", value)

    @property
    @pulumi.getter
    def verify_jwt(self) -> Optional[pulumi.Input[bool]]:
        """
        Verify JWT before running
        """
        return pulumi.get(self, "verify_jwt")

    @verify_jwt.setter
    def verify_jwt(self, value: Optional[pulumi.Input[bool]]):
        pulumi.set(self, "verify_jwt", value)


class Function(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 body: Optional[pulumi.Input[str]] = None,
                 name: Optional[pulumi.Input[str]] = None,
                 project_id: Optional[pulumi.Input[str]] = None,
                 slug: Optional[pulumi.Input[str]] = None,
                 verify_jwt: Optional[pulumi.Input[bool]] = None,
                 __props__=None):
        """
        Create a Function resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[str] body: Body of the functino
        :param pulumi.Input[str] name: Name of the function
        :param pulumi.Input[str] project_id: ID of the project
        :param pulumi.Input[str] slug: Slug of the function
        :param pulumi.Input[bool] verify_jwt: Verify JWT before running
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: FunctionArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a Function resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param FunctionArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(FunctionArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 body: Optional[pulumi.Input[str]] = None,
                 name: Optional[pulumi.Input[str]] = None,
                 project_id: Optional[pulumi.Input[str]] = None,
                 slug: Optional[pulumi.Input[str]] = None,
                 verify_jwt: Optional[pulumi.Input[bool]] = None,
                 __props__=None):
        if opts is None:
            opts = pulumi.ResourceOptions()
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.version is None:
            opts.version = _utilities.get_version()
        if opts.plugin_download_url is None:
            opts.plugin_download_url = _utilities.get_plugin_download_url()
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = FunctionArgs.__new__(FunctionArgs)

            if body is None and not opts.urn:
                raise TypeError("Missing required property 'body'")
            __props__.__dict__["body"] = None if body is None else pulumi.Output.secret(body)
            if name is None and not opts.urn:
                raise TypeError("Missing required property 'name'")
            __props__.__dict__["name"] = name
            if project_id is None and not opts.urn:
                raise TypeError("Missing required property 'project_id'")
            __props__.__dict__["project_id"] = project_id
            if slug is None and not opts.urn:
                raise TypeError("Missing required property 'slug'")
            __props__.__dict__["slug"] = slug
            if verify_jwt is None:
                verify_jwt = False
            __props__.__dict__["verify_jwt"] = verify_jwt
            __props__.__dict__["created_at"] = None
            __props__.__dict__["status"] = None
            __props__.__dict__["updated_at"] = None
            __props__.__dict__["version"] = None
        super(Function, __self__).__init__(
            'supabase:index:Function',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Function':
        """
        Get an existing Function resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = FunctionArgs.__new__(FunctionArgs)

        __props__.__dict__["created_at"] = None
        __props__.__dict__["name"] = None
        __props__.__dict__["slug"] = None
        __props__.__dict__["status"] = None
        __props__.__dict__["updated_at"] = None
        __props__.__dict__["verify_jwt"] = None
        __props__.__dict__["version"] = None
        return Function(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter
    def created_at(self) -> pulumi.Output[str]:
        """
        Function creation date
        """
        return pulumi.get(self, "created_at")

    @property
    @pulumi.getter
    def name(self) -> pulumi.Output[str]:
        """
        Name of the function
        """
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def slug(self) -> pulumi.Output[str]:
        """
        Slug of the function
        """
        return pulumi.get(self, "slug")

    @property
    @pulumi.getter
    def status(self) -> pulumi.Output['FunctionStatus']:
        """
        Status of the function
        """
        return pulumi.get(self, "status")

    @property
    @pulumi.getter(name="updatedAt")
    def updated_at(self) -> pulumi.Output[str]:
        """
        Function updated date
        """
        return pulumi.get(self, "updated_at")

    @property
    @pulumi.getter
    def verify_jwt(self) -> pulumi.Output[bool]:
        """
        Verify JWT before running
        """
        return pulumi.get(self, "verify_jwt")

    @property
    @pulumi.getter
    def version(self) -> pulumi.Output[int]:
        """
        Version of the function
        """
        return pulumi.get(self, "version")

