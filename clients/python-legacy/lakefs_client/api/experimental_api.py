"""
    lakeFS API

    lakeFS HTTP API  # noqa: E501

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by: https://openapi-generator.tech
"""


import re  # noqa: F401
import sys  # noqa: F401

from lakefs_client.api_client import ApiClient, Endpoint as _Endpoint
from lakefs_client.model_utils import (  # noqa: F401
    check_allowed_values,
    check_validations,
    date,
    datetime,
    file_type,
    none_type,
    validate_and_convert_types
)
from lakefs_client.model.abort_presign_multipart_upload import AbortPresignMultipartUpload
from lakefs_client.model.complete_presign_multipart_upload import CompletePresignMultipartUpload
from lakefs_client.model.error import Error
from lakefs_client.model.object_stats import ObjectStats
from lakefs_client.model.presign_multipart_upload import PresignMultipartUpload
from lakefs_client.model.staging_location import StagingLocation


class ExperimentalApi(object):
    """NOTE: This class is auto generated by OpenAPI Generator
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    def __init__(self, api_client=None):
        if api_client is None:
            api_client = ApiClient()
        self.api_client = api_client
        self.abort_presign_multipart_upload_endpoint = _Endpoint(
            settings={
                'response_type': None,
                'auth': [
                    'basic_auth',
                    'cookie_auth',
                    'jwt_token',
                    'oidc_auth',
                    'saml_auth'
                ],
                'endpoint_path': '/repositories/{repository}/branches/{branch}/staging/pmpu/{uploadId}',
                'operation_id': 'abort_presign_multipart_upload',
                'http_method': 'DELETE',
                'servers': None,
            },
            params_map={
                'all': [
                    'repository',
                    'branch',
                    'upload_id',
                    'path',
                    'abort_presign_multipart_upload',
                ],
                'required': [
                    'repository',
                    'branch',
                    'upload_id',
                    'path',
                ],
                'nullable': [
                ],
                'enum': [
                ],
                'validation': [
                ]
            },
            root_map={
                'validations': {
                },
                'allowed_values': {
                },
                'openapi_types': {
                    'repository':
                        (str,),
                    'branch':
                        (str,),
                    'upload_id':
                        (str,),
                    'path':
                        (str,),
                    'abort_presign_multipart_upload':
                        (AbortPresignMultipartUpload,),
                },
                'attribute_map': {
                    'repository': 'repository',
                    'branch': 'branch',
                    'upload_id': 'uploadId',
                    'path': 'path',
                },
                'location_map': {
                    'repository': 'path',
                    'branch': 'path',
                    'upload_id': 'path',
                    'path': 'query',
                    'abort_presign_multipart_upload': 'body',
                },
                'collection_format_map': {
                }
            },
            headers_map={
                'accept': [
                    'application/json'
                ],
                'content_type': [
                    'application/json'
                ]
            },
            api_client=api_client
        )
        self.complete_presign_multipart_upload_endpoint = _Endpoint(
            settings={
                'response_type': (ObjectStats,),
                'auth': [
                    'basic_auth',
                    'cookie_auth',
                    'jwt_token',
                    'oidc_auth',
                    'saml_auth'
                ],
                'endpoint_path': '/repositories/{repository}/branches/{branch}/staging/pmpu/{uploadId}',
                'operation_id': 'complete_presign_multipart_upload',
                'http_method': 'PUT',
                'servers': None,
            },
            params_map={
                'all': [
                    'repository',
                    'branch',
                    'upload_id',
                    'path',
                    'complete_presign_multipart_upload',
                ],
                'required': [
                    'repository',
                    'branch',
                    'upload_id',
                    'path',
                ],
                'nullable': [
                ],
                'enum': [
                ],
                'validation': [
                ]
            },
            root_map={
                'validations': {
                },
                'allowed_values': {
                },
                'openapi_types': {
                    'repository':
                        (str,),
                    'branch':
                        (str,),
                    'upload_id':
                        (str,),
                    'path':
                        (str,),
                    'complete_presign_multipart_upload':
                        (CompletePresignMultipartUpload,),
                },
                'attribute_map': {
                    'repository': 'repository',
                    'branch': 'branch',
                    'upload_id': 'uploadId',
                    'path': 'path',
                },
                'location_map': {
                    'repository': 'path',
                    'branch': 'path',
                    'upload_id': 'path',
                    'path': 'query',
                    'complete_presign_multipart_upload': 'body',
                },
                'collection_format_map': {
                }
            },
            headers_map={
                'accept': [
                    'application/json'
                ],
                'content_type': [
                    'application/json'
                ]
            },
            api_client=api_client
        )
        self.create_presign_multipart_upload_endpoint = _Endpoint(
            settings={
                'response_type': (PresignMultipartUpload,),
                'auth': [
                    'basic_auth',
                    'cookie_auth',
                    'jwt_token',
                    'oidc_auth',
                    'saml_auth'
                ],
                'endpoint_path': '/repositories/{repository}/branches/{branch}/staging/pmpu',
                'operation_id': 'create_presign_multipart_upload',
                'http_method': 'POST',
                'servers': None,
            },
            params_map={
                'all': [
                    'repository',
                    'branch',
                    'path',
                    'parts',
                ],
                'required': [
                    'repository',
                    'branch',
                    'path',
                ],
                'nullable': [
                ],
                'enum': [
                ],
                'validation': [
                ]
            },
            root_map={
                'validations': {
                },
                'allowed_values': {
                },
                'openapi_types': {
                    'repository':
                        (str,),
                    'branch':
                        (str,),
                    'path':
                        (str,),
                    'parts':
                        (int,),
                },
                'attribute_map': {
                    'repository': 'repository',
                    'branch': 'branch',
                    'path': 'path',
                    'parts': 'parts',
                },
                'location_map': {
                    'repository': 'path',
                    'branch': 'path',
                    'path': 'query',
                    'parts': 'query',
                },
                'collection_format_map': {
                }
            },
            headers_map={
                'accept': [
                    'application/json'
                ],
                'content_type': [],
            },
            api_client=api_client
        )
        self.hard_reset_branch_endpoint = _Endpoint(
            settings={
                'response_type': None,
                'auth': [
                    'basic_auth',
                    'cookie_auth',
                    'jwt_token',
                    'oidc_auth',
                    'saml_auth'
                ],
                'endpoint_path': '/repositories/{repository}/branches/{branch}/hard_reset',
                'operation_id': 'hard_reset_branch',
                'http_method': 'PUT',
                'servers': None,
            },
            params_map={
                'all': [
                    'repository',
                    'branch',
                    'ref',
                    'force',
                ],
                'required': [
                    'repository',
                    'branch',
                    'ref',
                ],
                'nullable': [
                ],
                'enum': [
                ],
                'validation': [
                ]
            },
            root_map={
                'validations': {
                },
                'allowed_values': {
                },
                'openapi_types': {
                    'repository':
                        (str,),
                    'branch':
                        (str,),
                    'ref':
                        (str,),
                    'force':
                        (bool,),
                },
                'attribute_map': {
                    'repository': 'repository',
                    'branch': 'branch',
                    'ref': 'ref',
                    'force': 'force',
                },
                'location_map': {
                    'repository': 'path',
                    'branch': 'path',
                    'ref': 'query',
                    'force': 'query',
                },
                'collection_format_map': {
                }
            },
            headers_map={
                'accept': [
                    'application/json'
                ],
                'content_type': [],
            },
            api_client=api_client
        )

    def abort_presign_multipart_upload(
        self,
        repository,
        branch,
        upload_id,
        path,
        **kwargs
    ):
        """Abort a presign multipart upload  # noqa: E501

        Aborts a presign multipart upload.  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.abort_presign_multipart_upload(repository, branch, upload_id, path, async_req=True)
        >>> result = thread.get()

        Args:
            repository (str):
            branch (str):
            upload_id (str):
            path (str): relative to the branch

        Keyword Args:
            abort_presign_multipart_upload (AbortPresignMultipartUpload): [optional]
            _return_http_data_only (bool): response data without head status
                code and headers. Default is True.
            _preload_content (bool): if False, the urllib3.HTTPResponse object
                will be returned without reading/decoding response data.
                Default is True.
            _request_timeout (int/float/tuple): timeout setting for this request. If
                one number provided, it will be total request timeout. It can also
                be a pair (tuple) of (connection, read) timeouts.
                Default is None.
            _check_input_type (bool): specifies if type checking
                should be done one the data sent to the server.
                Default is True.
            _check_return_type (bool): specifies if type checking
                should be done one the data received from the server.
                Default is True.
            _host_index (int/None): specifies the index of the server
                that we want to use.
                Default is read from the configuration.
            async_req (bool): execute request asynchronously

        Returns:
            None
                If the method is called asynchronously, returns the request
                thread.
        """
        kwargs['async_req'] = kwargs.get(
            'async_req', False
        )
        kwargs['_return_http_data_only'] = kwargs.get(
            '_return_http_data_only', True
        )
        kwargs['_preload_content'] = kwargs.get(
            '_preload_content', True
        )
        kwargs['_request_timeout'] = kwargs.get(
            '_request_timeout', None
        )
        kwargs['_check_input_type'] = kwargs.get(
            '_check_input_type', True
        )
        kwargs['_check_return_type'] = kwargs.get(
            '_check_return_type', True
        )
        kwargs['_host_index'] = kwargs.get('_host_index')
        kwargs['repository'] = \
            repository
        kwargs['branch'] = \
            branch
        kwargs['upload_id'] = \
            upload_id
        kwargs['path'] = \
            path
        return self.abort_presign_multipart_upload_endpoint.call_with_http_info(**kwargs)

    def complete_presign_multipart_upload(
        self,
        repository,
        branch,
        upload_id,
        path,
        **kwargs
    ):
        """Complete a presign multipart upload request  # noqa: E501

        Completes a presign multipart upload by assembling the uploaded parts.  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.complete_presign_multipart_upload(repository, branch, upload_id, path, async_req=True)
        >>> result = thread.get()

        Args:
            repository (str):
            branch (str):
            upload_id (str):
            path (str): relative to the branch

        Keyword Args:
            complete_presign_multipart_upload (CompletePresignMultipartUpload): [optional]
            _return_http_data_only (bool): response data without head status
                code and headers. Default is True.
            _preload_content (bool): if False, the urllib3.HTTPResponse object
                will be returned without reading/decoding response data.
                Default is True.
            _request_timeout (int/float/tuple): timeout setting for this request. If
                one number provided, it will be total request timeout. It can also
                be a pair (tuple) of (connection, read) timeouts.
                Default is None.
            _check_input_type (bool): specifies if type checking
                should be done one the data sent to the server.
                Default is True.
            _check_return_type (bool): specifies if type checking
                should be done one the data received from the server.
                Default is True.
            _host_index (int/None): specifies the index of the server
                that we want to use.
                Default is read from the configuration.
            async_req (bool): execute request asynchronously

        Returns:
            ObjectStats
                If the method is called asynchronously, returns the request
                thread.
        """
        kwargs['async_req'] = kwargs.get(
            'async_req', False
        )
        kwargs['_return_http_data_only'] = kwargs.get(
            '_return_http_data_only', True
        )
        kwargs['_preload_content'] = kwargs.get(
            '_preload_content', True
        )
        kwargs['_request_timeout'] = kwargs.get(
            '_request_timeout', None
        )
        kwargs['_check_input_type'] = kwargs.get(
            '_check_input_type', True
        )
        kwargs['_check_return_type'] = kwargs.get(
            '_check_return_type', True
        )
        kwargs['_host_index'] = kwargs.get('_host_index')
        kwargs['repository'] = \
            repository
        kwargs['branch'] = \
            branch
        kwargs['upload_id'] = \
            upload_id
        kwargs['path'] = \
            path
        return self.complete_presign_multipart_upload_endpoint.call_with_http_info(**kwargs)

    def create_presign_multipart_upload(
        self,
        repository,
        branch,
        path,
        **kwargs
    ):
        """Initiate a multipart upload  # noqa: E501

        Initiates a multipart upload and returns an upload ID with presigned URLs for each part (optional). Part numbers starts with 1. Each part except the last one has minimum size depends on the underlying blockstore implementation. For example working with S3 blockstore, minimum size is 5MB (excluding the last part).   # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.create_presign_multipart_upload(repository, branch, path, async_req=True)
        >>> result = thread.get()

        Args:
            repository (str):
            branch (str):
            path (str): relative to the branch

        Keyword Args:
            parts (int): number of presigned URL parts required to upload. [optional]
            _return_http_data_only (bool): response data without head status
                code and headers. Default is True.
            _preload_content (bool): if False, the urllib3.HTTPResponse object
                will be returned without reading/decoding response data.
                Default is True.
            _request_timeout (int/float/tuple): timeout setting for this request. If
                one number provided, it will be total request timeout. It can also
                be a pair (tuple) of (connection, read) timeouts.
                Default is None.
            _check_input_type (bool): specifies if type checking
                should be done one the data sent to the server.
                Default is True.
            _check_return_type (bool): specifies if type checking
                should be done one the data received from the server.
                Default is True.
            _host_index (int/None): specifies the index of the server
                that we want to use.
                Default is read from the configuration.
            async_req (bool): execute request asynchronously

        Returns:
            PresignMultipartUpload
                If the method is called asynchronously, returns the request
                thread.
        """
        kwargs['async_req'] = kwargs.get(
            'async_req', False
        )
        kwargs['_return_http_data_only'] = kwargs.get(
            '_return_http_data_only', True
        )
        kwargs['_preload_content'] = kwargs.get(
            '_preload_content', True
        )
        kwargs['_request_timeout'] = kwargs.get(
            '_request_timeout', None
        )
        kwargs['_check_input_type'] = kwargs.get(
            '_check_input_type', True
        )
        kwargs['_check_return_type'] = kwargs.get(
            '_check_return_type', True
        )
        kwargs['_host_index'] = kwargs.get('_host_index')
        kwargs['repository'] = \
            repository
        kwargs['branch'] = \
            branch
        kwargs['path'] = \
            path
        return self.create_presign_multipart_upload_endpoint.call_with_http_info(**kwargs)

    def hard_reset_branch(
        self,
        repository,
        branch,
        ref,
        **kwargs
    ):
        """hard reset branch  # noqa: E501

        Relocate branch to refer to ref.  Branch must not contain uncommitted data.  # noqa: E501
        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.hard_reset_branch(repository, branch, ref, async_req=True)
        >>> result = thread.get()

        Args:
            repository (str):
            branch (str):
            ref (str): After reset, branch will point at this reference.

        Keyword Args:
            force (bool): [optional] if omitted the server will use the default value of False
            _return_http_data_only (bool): response data without head status
                code and headers. Default is True.
            _preload_content (bool): if False, the urllib3.HTTPResponse object
                will be returned without reading/decoding response data.
                Default is True.
            _request_timeout (int/float/tuple): timeout setting for this request. If
                one number provided, it will be total request timeout. It can also
                be a pair (tuple) of (connection, read) timeouts.
                Default is None.
            _check_input_type (bool): specifies if type checking
                should be done one the data sent to the server.
                Default is True.
            _check_return_type (bool): specifies if type checking
                should be done one the data received from the server.
                Default is True.
            _host_index (int/None): specifies the index of the server
                that we want to use.
                Default is read from the configuration.
            async_req (bool): execute request asynchronously

        Returns:
            None
                If the method is called asynchronously, returns the request
                thread.
        """
        kwargs['async_req'] = kwargs.get(
            'async_req', False
        )
        kwargs['_return_http_data_only'] = kwargs.get(
            '_return_http_data_only', True
        )
        kwargs['_preload_content'] = kwargs.get(
            '_preload_content', True
        )
        kwargs['_request_timeout'] = kwargs.get(
            '_request_timeout', None
        )
        kwargs['_check_input_type'] = kwargs.get(
            '_check_input_type', True
        )
        kwargs['_check_return_type'] = kwargs.get(
            '_check_return_type', True
        )
        kwargs['_host_index'] = kwargs.get('_host_index')
        kwargs['repository'] = \
            repository
        kwargs['branch'] = \
            branch
        kwargs['ref'] = \
            ref
        return self.hard_reset_branch_endpoint.call_with_http_info(**kwargs)

