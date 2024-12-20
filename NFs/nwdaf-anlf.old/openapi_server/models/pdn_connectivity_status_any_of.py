# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class PdnConnectivityStatusAnyOf(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    """
    allowed enum values
    """
    ESTABLISHED = "ESTABLISHED"
    RELEASED = "RELEASED"
    def __init__(self):  # noqa: E501
        """PdnConnectivityStatusAnyOf - a model defined in OpenAPI

        """
        self.openapi_types = {
        }

        self.attribute_map = {
        }

    @classmethod
    def from_dict(cls, dikt) -> 'PdnConnectivityStatusAnyOf':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PdnConnectivityStatus_anyOf of this PdnConnectivityStatusAnyOf.  # noqa: E501
        :rtype: PdnConnectivityStatusAnyOf
        """
        return util.deserialize_model(dikt, cls)
