# coding: utf-8

"""
    lakeFS API

    lakeFS HTTP API

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import re  # noqa: F401
import io
import warnings

try:
    from pydantic.v1 import validate_arguments, ValidationError
except ImportError:
    from pydantic import validate_arguments, ValidationError
from typing_extensions import Annotated

from datetime import datetime

try:
    from pydantic.v1 import Field, StrictBool, StrictStr, conint, conlist
except ImportError:
    from pydantic import Field, StrictBool, StrictStr, conint, conlist

from typing import Optional

from lakefs_sdk.models.commit_list import CommitList
from lakefs_sdk.models.diff_list import DiffList
from lakefs_sdk.models.find_merge_base_result import FindMergeBaseResult
from lakefs_sdk.models.merge import Merge
from lakefs_sdk.models.merge_result import MergeResult

from lakefs_sdk.api_client import ApiClient
from lakefs_sdk.api_response import ApiResponse
from lakefs_sdk.exceptions import (  # noqa: F401
    ApiTypeError,
    ApiValueError
)


class RefsApi(object):
    """NOTE: This class is auto generated by OpenAPI Generator
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    def __init__(self, api_client=None):
        if api_client is None:
            api_client = ApiClient.get_default()
        self.api_client = api_client

    @validate_arguments
    def diff_refs(self, repository : StrictStr, left_ref : Annotated[StrictStr, Field(..., description="a reference (could be either a branch or a commit ID)")], right_ref : Annotated[StrictStr, Field(..., description="a reference (could be either a branch or a commit ID) to compare against")], after : Annotated[Optional[StrictStr], Field(description="return items after this value")] = None, amount : Annotated[Optional[conint(strict=True, le=1000, ge=-1)], Field(description="how many items to return")] = None, prefix : Annotated[Optional[StrictStr], Field(description="return items prefixed with this value")] = None, delimiter : Annotated[Optional[StrictStr], Field(description="delimiter used to group common prefixes by")] = None, type : Optional[StrictStr] = None, **kwargs) -> DiffList:  # noqa: E501
        """diff references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.diff_refs(repository, left_ref, right_ref, after, amount, prefix, delimiter, type, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param left_ref: a reference (could be either a branch or a commit ID) (required)
        :type left_ref: str
        :param right_ref: a reference (could be either a branch or a commit ID) to compare against (required)
        :type right_ref: str
        :param after: return items after this value
        :type after: str
        :param amount: how many items to return
        :type amount: int
        :param prefix: return items prefixed with this value
        :type prefix: str
        :param delimiter: delimiter used to group common prefixes by
        :type delimiter: str
        :param type:
        :type type: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: DiffList
        """
        kwargs['_return_http_data_only'] = True
        if '_preload_content' in kwargs:
            raise ValueError("Error! Please call the diff_refs_with_http_info method with `_preload_content` instead and obtain raw data from ApiResponse.raw_data")
        return self.diff_refs_with_http_info(repository, left_ref, right_ref, after, amount, prefix, delimiter, type, **kwargs)  # noqa: E501

    @validate_arguments
    def diff_refs_with_http_info(self, repository : StrictStr, left_ref : Annotated[StrictStr, Field(..., description="a reference (could be either a branch or a commit ID)")], right_ref : Annotated[StrictStr, Field(..., description="a reference (could be either a branch or a commit ID) to compare against")], after : Annotated[Optional[StrictStr], Field(description="return items after this value")] = None, amount : Annotated[Optional[conint(strict=True, le=1000, ge=-1)], Field(description="how many items to return")] = None, prefix : Annotated[Optional[StrictStr], Field(description="return items prefixed with this value")] = None, delimiter : Annotated[Optional[StrictStr], Field(description="delimiter used to group common prefixes by")] = None, type : Optional[StrictStr] = None, **kwargs) -> ApiResponse:  # noqa: E501
        """diff references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.diff_refs_with_http_info(repository, left_ref, right_ref, after, amount, prefix, delimiter, type, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param left_ref: a reference (could be either a branch or a commit ID) (required)
        :type left_ref: str
        :param right_ref: a reference (could be either a branch or a commit ID) to compare against (required)
        :type right_ref: str
        :param after: return items after this value
        :type after: str
        :param amount: how many items to return
        :type amount: int
        :param prefix: return items prefixed with this value
        :type prefix: str
        :param delimiter: delimiter used to group common prefixes by
        :type delimiter: str
        :param type:
        :type type: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _preload_content: if False, the ApiResponse.data will
                                 be set to none and raw_data will store the 
                                 HTTP response body without reading/decoding.
                                 Default is True.
        :type _preload_content: bool, optional
        :param _return_http_data_only: response data instead of ApiResponse
                                       object with status code, headers, etc
        :type _return_http_data_only: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :param _request_auth: set to override the auth_settings for an a single
                              request; this effectively ignores the authentication
                              in the spec for a single request.
        :type _request_auth: dict, optional
        :type _content_type: string, optional: force content-type for the request
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: tuple(DiffList, status_code(int), headers(HTTPHeaderDict))
        """

        _params = locals()

        _all_params = [
            'repository',
            'left_ref',
            'right_ref',
            'after',
            'amount',
            'prefix',
            'delimiter',
            'type'
        ]
        _all_params.extend(
            [
                'async_req',
                '_return_http_data_only',
                '_preload_content',
                '_request_timeout',
                '_request_auth',
                '_content_type',
                '_headers'
            ]
        )

        # validate the arguments
        for _key, _val in _params['kwargs'].items():
            if _key not in _all_params:
                raise ApiTypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method diff_refs" % _key
                )
            _params[_key] = _val
        del _params['kwargs']

        _collection_formats = {}

        # process the path parameters
        _path_params = {}
        if _params['repository']:
            _path_params['repository'] = _params['repository']

        if _params['left_ref']:
            _path_params['leftRef'] = _params['left_ref']

        if _params['right_ref']:
            _path_params['rightRef'] = _params['right_ref']


        # process the query parameters
        _query_params = []
        if _params.get('after') is not None:  # noqa: E501
            _query_params.append(('after', _params['after']))

        if _params.get('amount') is not None:  # noqa: E501
            _query_params.append(('amount', _params['amount']))

        if _params.get('prefix') is not None:  # noqa: E501
            _query_params.append(('prefix', _params['prefix']))

        if _params.get('delimiter') is not None:  # noqa: E501
            _query_params.append(('delimiter', _params['delimiter']))

        if _params.get('type') is not None:  # noqa: E501
            _query_params.append(('type', _params['type']))

        # process the header parameters
        _header_params = dict(_params.get('_headers', {}))
        # process the form parameters
        _form_params = []
        _files = {}
        # process the body parameter
        _body_params = None
        # set the HTTP header `Accept`
        _header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # authentication setting
        _auth_settings = ['basic_auth', 'cookie_auth', 'oidc_auth', 'saml_auth', 'jwt_token']  # noqa: E501

        _response_types_map = {
            '200': "DiffList",
            '401': "Error",
            '404': "Error",
            '420': None,
        }

        return self.api_client.call_api(
            '/repositories/{repository}/refs/{leftRef}/diff/{rightRef}', 'GET',
            _path_params,
            _query_params,
            _header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            response_types_map=_response_types_map,
            auth_settings=_auth_settings,
            async_req=_params.get('async_req'),
            _return_http_data_only=_params.get('_return_http_data_only'),  # noqa: E501
            _preload_content=_params.get('_preload_content', True),
            _request_timeout=_params.get('_request_timeout'),
            collection_formats=_collection_formats,
            _request_auth=_params.get('_request_auth'))

    @validate_arguments
    def find_merge_base(self, repository : StrictStr, source_ref : Annotated[StrictStr, Field(..., description="source ref")], destination_branch : Annotated[StrictStr, Field(..., description="destination branch name")], **kwargs) -> FindMergeBaseResult:  # noqa: E501
        """find the merge base for 2 references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.find_merge_base(repository, source_ref, destination_branch, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param source_ref: source ref (required)
        :type source_ref: str
        :param destination_branch: destination branch name (required)
        :type destination_branch: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: FindMergeBaseResult
        """
        kwargs['_return_http_data_only'] = True
        if '_preload_content' in kwargs:
            raise ValueError("Error! Please call the find_merge_base_with_http_info method with `_preload_content` instead and obtain raw data from ApiResponse.raw_data")
        return self.find_merge_base_with_http_info(repository, source_ref, destination_branch, **kwargs)  # noqa: E501

    @validate_arguments
    def find_merge_base_with_http_info(self, repository : StrictStr, source_ref : Annotated[StrictStr, Field(..., description="source ref")], destination_branch : Annotated[StrictStr, Field(..., description="destination branch name")], **kwargs) -> ApiResponse:  # noqa: E501
        """find the merge base for 2 references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.find_merge_base_with_http_info(repository, source_ref, destination_branch, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param source_ref: source ref (required)
        :type source_ref: str
        :param destination_branch: destination branch name (required)
        :type destination_branch: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _preload_content: if False, the ApiResponse.data will
                                 be set to none and raw_data will store the 
                                 HTTP response body without reading/decoding.
                                 Default is True.
        :type _preload_content: bool, optional
        :param _return_http_data_only: response data instead of ApiResponse
                                       object with status code, headers, etc
        :type _return_http_data_only: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :param _request_auth: set to override the auth_settings for an a single
                              request; this effectively ignores the authentication
                              in the spec for a single request.
        :type _request_auth: dict, optional
        :type _content_type: string, optional: force content-type for the request
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: tuple(FindMergeBaseResult, status_code(int), headers(HTTPHeaderDict))
        """

        _params = locals()

        _all_params = [
            'repository',
            'source_ref',
            'destination_branch'
        ]
        _all_params.extend(
            [
                'async_req',
                '_return_http_data_only',
                '_preload_content',
                '_request_timeout',
                '_request_auth',
                '_content_type',
                '_headers'
            ]
        )

        # validate the arguments
        for _key, _val in _params['kwargs'].items():
            if _key not in _all_params:
                raise ApiTypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method find_merge_base" % _key
                )
            _params[_key] = _val
        del _params['kwargs']

        _collection_formats = {}

        # process the path parameters
        _path_params = {}
        if _params['repository']:
            _path_params['repository'] = _params['repository']

        if _params['source_ref']:
            _path_params['sourceRef'] = _params['source_ref']

        if _params['destination_branch']:
            _path_params['destinationBranch'] = _params['destination_branch']


        # process the query parameters
        _query_params = []
        # process the header parameters
        _header_params = dict(_params.get('_headers', {}))
        # process the form parameters
        _form_params = []
        _files = {}
        # process the body parameter
        _body_params = None
        # set the HTTP header `Accept`
        _header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # authentication setting
        _auth_settings = ['basic_auth', 'cookie_auth', 'oidc_auth', 'saml_auth', 'jwt_token']  # noqa: E501

        _response_types_map = {
            '200': "FindMergeBaseResult",
            '400': "Error",
            '401': "Error",
            '404': "Error",
            '420': None,
        }

        return self.api_client.call_api(
            '/repositories/{repository}/refs/{sourceRef}/merge/{destinationBranch}', 'GET',
            _path_params,
            _query_params,
            _header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            response_types_map=_response_types_map,
            auth_settings=_auth_settings,
            async_req=_params.get('async_req'),
            _return_http_data_only=_params.get('_return_http_data_only'),  # noqa: E501
            _preload_content=_params.get('_preload_content', True),
            _request_timeout=_params.get('_request_timeout'),
            collection_formats=_collection_formats,
            _request_auth=_params.get('_request_auth'))

    @validate_arguments
    def log_commits(self, repository : StrictStr, ref : StrictStr, after : Annotated[Optional[StrictStr], Field(description="return items after this value")] = None, amount : Annotated[Optional[conint(strict=True, le=1000, ge=-1)], Field(description="how many items to return")] = None, objects : Annotated[Optional[conlist(StrictStr)], Field(description="list of paths, each element is a path of a specific object")] = None, prefixes : Annotated[Optional[conlist(StrictStr)], Field(description="list of paths, each element is a path of a prefix")] = None, limit : Annotated[Optional[StrictBool], Field(description="limit the number of items in return to 'amount'. Without further indication on actual number of items.")] = None, first_parent : Annotated[Optional[StrictBool], Field(description="if set to true, follow only the first parent upon reaching a merge commit")] = None, since : Annotated[Optional[datetime], Field(description="Show commits more recent than a specific date-time. In case used with stop_at parameter, will stop at the first commit that meets any of the conditions.")] = None, stop_at : Annotated[Optional[StrictStr], Field(description="A reference to stop at. In case used with since parameter, will stop at the first commit that meets any of the conditions.")] = None, **kwargs) -> CommitList:  # noqa: E501
        """get commit log from ref. If both objects and prefixes are empty, return all commits.  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.log_commits(repository, ref, after, amount, objects, prefixes, limit, first_parent, since, stop_at, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param ref: (required)
        :type ref: str
        :param after: return items after this value
        :type after: str
        :param amount: how many items to return
        :type amount: int
        :param objects: list of paths, each element is a path of a specific object
        :type objects: List[str]
        :param prefixes: list of paths, each element is a path of a prefix
        :type prefixes: List[str]
        :param limit: limit the number of items in return to 'amount'. Without further indication on actual number of items.
        :type limit: bool
        :param first_parent: if set to true, follow only the first parent upon reaching a merge commit
        :type first_parent: bool
        :param since: Show commits more recent than a specific date-time. In case used with stop_at parameter, will stop at the first commit that meets any of the conditions.
        :type since: datetime
        :param stop_at: A reference to stop at. In case used with since parameter, will stop at the first commit that meets any of the conditions.
        :type stop_at: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: CommitList
        """
        kwargs['_return_http_data_only'] = True
        if '_preload_content' in kwargs:
            raise ValueError("Error! Please call the log_commits_with_http_info method with `_preload_content` instead and obtain raw data from ApiResponse.raw_data")
        return self.log_commits_with_http_info(repository, ref, after, amount, objects, prefixes, limit, first_parent, since, stop_at, **kwargs)  # noqa: E501

    @validate_arguments
    def log_commits_with_http_info(self, repository : StrictStr, ref : StrictStr, after : Annotated[Optional[StrictStr], Field(description="return items after this value")] = None, amount : Annotated[Optional[conint(strict=True, le=1000, ge=-1)], Field(description="how many items to return")] = None, objects : Annotated[Optional[conlist(StrictStr)], Field(description="list of paths, each element is a path of a specific object")] = None, prefixes : Annotated[Optional[conlist(StrictStr)], Field(description="list of paths, each element is a path of a prefix")] = None, limit : Annotated[Optional[StrictBool], Field(description="limit the number of items in return to 'amount'. Without further indication on actual number of items.")] = None, first_parent : Annotated[Optional[StrictBool], Field(description="if set to true, follow only the first parent upon reaching a merge commit")] = None, since : Annotated[Optional[datetime], Field(description="Show commits more recent than a specific date-time. In case used with stop_at parameter, will stop at the first commit that meets any of the conditions.")] = None, stop_at : Annotated[Optional[StrictStr], Field(description="A reference to stop at. In case used with since parameter, will stop at the first commit that meets any of the conditions.")] = None, **kwargs) -> ApiResponse:  # noqa: E501
        """get commit log from ref. If both objects and prefixes are empty, return all commits.  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.log_commits_with_http_info(repository, ref, after, amount, objects, prefixes, limit, first_parent, since, stop_at, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param ref: (required)
        :type ref: str
        :param after: return items after this value
        :type after: str
        :param amount: how many items to return
        :type amount: int
        :param objects: list of paths, each element is a path of a specific object
        :type objects: List[str]
        :param prefixes: list of paths, each element is a path of a prefix
        :type prefixes: List[str]
        :param limit: limit the number of items in return to 'amount'. Without further indication on actual number of items.
        :type limit: bool
        :param first_parent: if set to true, follow only the first parent upon reaching a merge commit
        :type first_parent: bool
        :param since: Show commits more recent than a specific date-time. In case used with stop_at parameter, will stop at the first commit that meets any of the conditions.
        :type since: datetime
        :param stop_at: A reference to stop at. In case used with since parameter, will stop at the first commit that meets any of the conditions.
        :type stop_at: str
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _preload_content: if False, the ApiResponse.data will
                                 be set to none and raw_data will store the 
                                 HTTP response body without reading/decoding.
                                 Default is True.
        :type _preload_content: bool, optional
        :param _return_http_data_only: response data instead of ApiResponse
                                       object with status code, headers, etc
        :type _return_http_data_only: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :param _request_auth: set to override the auth_settings for an a single
                              request; this effectively ignores the authentication
                              in the spec for a single request.
        :type _request_auth: dict, optional
        :type _content_type: string, optional: force content-type for the request
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: tuple(CommitList, status_code(int), headers(HTTPHeaderDict))
        """

        _params = locals()

        _all_params = [
            'repository',
            'ref',
            'after',
            'amount',
            'objects',
            'prefixes',
            'limit',
            'first_parent',
            'since',
            'stop_at'
        ]
        _all_params.extend(
            [
                'async_req',
                '_return_http_data_only',
                '_preload_content',
                '_request_timeout',
                '_request_auth',
                '_content_type',
                '_headers'
            ]
        )

        # validate the arguments
        for _key, _val in _params['kwargs'].items():
            if _key not in _all_params:
                raise ApiTypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method log_commits" % _key
                )
            _params[_key] = _val
        del _params['kwargs']

        _collection_formats = {}

        # process the path parameters
        _path_params = {}
        if _params['repository']:
            _path_params['repository'] = _params['repository']

        if _params['ref']:
            _path_params['ref'] = _params['ref']


        # process the query parameters
        _query_params = []
        if _params.get('after') is not None:  # noqa: E501
            _query_params.append(('after', _params['after']))

        if _params.get('amount') is not None:  # noqa: E501
            _query_params.append(('amount', _params['amount']))

        if _params.get('objects') is not None:  # noqa: E501
            _query_params.append(('objects', _params['objects']))
            _collection_formats['objects'] = 'multi'

        if _params.get('prefixes') is not None:  # noqa: E501
            _query_params.append(('prefixes', _params['prefixes']))
            _collection_formats['prefixes'] = 'multi'

        if _params.get('limit') is not None:  # noqa: E501
            _query_params.append(('limit', _params['limit']))

        if _params.get('first_parent') is not None:  # noqa: E501
            _query_params.append(('first_parent', _params['first_parent']))

        if _params.get('since') is not None:  # noqa: E501
            if isinstance(_params['since'], datetime):
                _query_params.append(('since', _params['since'].strftime(self.api_client.configuration.datetime_format)))
            else:
                _query_params.append(('since', _params['since']))

        if _params.get('stop_at') is not None:  # noqa: E501
            _query_params.append(('stop_at', _params['stop_at']))

        # process the header parameters
        _header_params = dict(_params.get('_headers', {}))
        # process the form parameters
        _form_params = []
        _files = {}
        # process the body parameter
        _body_params = None
        # set the HTTP header `Accept`
        _header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # authentication setting
        _auth_settings = ['basic_auth', 'cookie_auth', 'oidc_auth', 'saml_auth', 'jwt_token']  # noqa: E501

        _response_types_map = {
            '200': "CommitList",
            '401': "Error",
            '404': "Error",
            '420': None,
        }

        return self.api_client.call_api(
            '/repositories/{repository}/refs/{ref}/commits', 'GET',
            _path_params,
            _query_params,
            _header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            response_types_map=_response_types_map,
            auth_settings=_auth_settings,
            async_req=_params.get('async_req'),
            _return_http_data_only=_params.get('_return_http_data_only'),  # noqa: E501
            _preload_content=_params.get('_preload_content', True),
            _request_timeout=_params.get('_request_timeout'),
            collection_formats=_collection_formats,
            _request_auth=_params.get('_request_auth'))

    @validate_arguments
    def merge_into_branch(self, repository : StrictStr, source_ref : Annotated[StrictStr, Field(..., description="source ref")], destination_branch : Annotated[StrictStr, Field(..., description="destination branch name")], merge : Optional[Merge] = None, **kwargs) -> MergeResult:  # noqa: E501
        """merge references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.merge_into_branch(repository, source_ref, destination_branch, merge, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param source_ref: source ref (required)
        :type source_ref: str
        :param destination_branch: destination branch name (required)
        :type destination_branch: str
        :param merge:
        :type merge: Merge
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: MergeResult
        """
        kwargs['_return_http_data_only'] = True
        if '_preload_content' in kwargs:
            raise ValueError("Error! Please call the merge_into_branch_with_http_info method with `_preload_content` instead and obtain raw data from ApiResponse.raw_data")
        return self.merge_into_branch_with_http_info(repository, source_ref, destination_branch, merge, **kwargs)  # noqa: E501

    @validate_arguments
    def merge_into_branch_with_http_info(self, repository : StrictStr, source_ref : Annotated[StrictStr, Field(..., description="source ref")], destination_branch : Annotated[StrictStr, Field(..., description="destination branch name")], merge : Optional[Merge] = None, **kwargs) -> ApiResponse:  # noqa: E501
        """merge references  # noqa: E501

        This method makes a synchronous HTTP request by default. To make an
        asynchronous HTTP request, please pass async_req=True

        >>> thread = api.merge_into_branch_with_http_info(repository, source_ref, destination_branch, merge, async_req=True)
        >>> result = thread.get()

        :param repository: (required)
        :type repository: str
        :param source_ref: source ref (required)
        :type source_ref: str
        :param destination_branch: destination branch name (required)
        :type destination_branch: str
        :param merge:
        :type merge: Merge
        :param async_req: Whether to execute the request asynchronously.
        :type async_req: bool, optional
        :param _preload_content: if False, the ApiResponse.data will
                                 be set to none and raw_data will store the 
                                 HTTP response body without reading/decoding.
                                 Default is True.
        :type _preload_content: bool, optional
        :param _return_http_data_only: response data instead of ApiResponse
                                       object with status code, headers, etc
        :type _return_http_data_only: bool, optional
        :param _request_timeout: timeout setting for this request. If one
                                 number provided, it will be total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :param _request_auth: set to override the auth_settings for an a single
                              request; this effectively ignores the authentication
                              in the spec for a single request.
        :type _request_auth: dict, optional
        :type _content_type: string, optional: force content-type for the request
        :return: Returns the result object.
                 If the method is called asynchronously,
                 returns the request thread.
        :rtype: tuple(MergeResult, status_code(int), headers(HTTPHeaderDict))
        """

        _params = locals()

        _all_params = [
            'repository',
            'source_ref',
            'destination_branch',
            'merge'
        ]
        _all_params.extend(
            [
                'async_req',
                '_return_http_data_only',
                '_preload_content',
                '_request_timeout',
                '_request_auth',
                '_content_type',
                '_headers'
            ]
        )

        # validate the arguments
        for _key, _val in _params['kwargs'].items():
            if _key not in _all_params:
                raise ApiTypeError(
                    "Got an unexpected keyword argument '%s'"
                    " to method merge_into_branch" % _key
                )
            _params[_key] = _val
        del _params['kwargs']

        _collection_formats = {}

        # process the path parameters
        _path_params = {}
        if _params['repository']:
            _path_params['repository'] = _params['repository']

        if _params['source_ref']:
            _path_params['sourceRef'] = _params['source_ref']

        if _params['destination_branch']:
            _path_params['destinationBranch'] = _params['destination_branch']


        # process the query parameters
        _query_params = []
        # process the header parameters
        _header_params = dict(_params.get('_headers', {}))
        # process the form parameters
        _form_params = []
        _files = {}
        # process the body parameter
        _body_params = None
        if _params['merge'] is not None:
            _body_params = _params['merge']

        # set the HTTP header `Accept`
        _header_params['Accept'] = self.api_client.select_header_accept(
            ['application/json'])  # noqa: E501

        # set the HTTP header `Content-Type`
        _content_types_list = _params.get('_content_type',
            self.api_client.select_header_content_type(
                ['application/json']))
        if _content_types_list:
                _header_params['Content-Type'] = _content_types_list

        # authentication setting
        _auth_settings = ['basic_auth', 'cookie_auth', 'oidc_auth', 'saml_auth', 'jwt_token']  # noqa: E501

        _response_types_map = {
            '200': "MergeResult",
            '400': "Error",
            '401': "Error",
            '403': "Error",
            '404': "Error",
            '409': "MergeResult",
            '412': "Error",
            '420': None,
        }

        return self.api_client.call_api(
            '/repositories/{repository}/refs/{sourceRef}/merge/{destinationBranch}', 'POST',
            _path_params,
            _query_params,
            _header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            response_types_map=_response_types_map,
            auth_settings=_auth_settings,
            async_req=_params.get('async_req'),
            _return_http_data_only=_params.get('_return_http_data_only'),  # noqa: E501
            _preload_content=_params.get('_preload_content', True),
            _request_timeout=_params.get('_request_timeout'),
            collection_formats=_collection_formats,
            _request_auth=_params.get('_request_auth'))
