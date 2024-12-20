# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.rm_state_any_of import RmStateAnyOf
from openapi_server import util

from openapi_server.models.rm_state_any_of import RmStateAnyOf  # noqa: E501

class RmState(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self):  # noqa: E501
        """RmState - a model defined in OpenAPI

        """
        self.openapi_types = {
        }

        self.attribute_map = {
        }

    @classmethod
    def from_dict(cls, dikt) -> 'RmState':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The RmState of this RmState.  # noqa: E501
        :rtype: RmState
        """
        return util.deserialize_model(dikt, cls)
