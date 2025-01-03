# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ue_type import UeType
from openapi_server import util

from openapi_server.models.ue_type import UeType  # noqa: E501

class UeInAreaFilter(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ue_type=None, aerial_srv_dnn_ind=False):  # noqa: E501
        """UeInAreaFilter - a model defined in OpenAPI

        :param ue_type: The ue_type of this UeInAreaFilter.  # noqa: E501
        :type ue_type: UeType
        :param aerial_srv_dnn_ind: The aerial_srv_dnn_ind of this UeInAreaFilter.  # noqa: E501
        :type aerial_srv_dnn_ind: bool
        """
        self.openapi_types = {
            'ue_type': UeType,
            'aerial_srv_dnn_ind': bool
        }

        self.attribute_map = {
            'ue_type': 'ueType',
            'aerial_srv_dnn_ind': 'aerialSrvDnnInd'
        }

        self.ue_type = ue_type
        self.aerial_srv_dnn_ind = aerial_srv_dnn_ind

    @classmethod
    def from_dict(cls, dikt) -> 'UeInAreaFilter':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The UeInAreaFilter of this UeInAreaFilter.  # noqa: E501
        :rtype: UeInAreaFilter
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ue_type(self):
        """Gets the ue_type of this UeInAreaFilter.


        :return: The ue_type of this UeInAreaFilter.
        :rtype: UeType
        """
        return self._ue_type

    @ue_type.setter
    def ue_type(self, ue_type):
        """Sets the ue_type of this UeInAreaFilter.


        :param ue_type: The ue_type of this UeInAreaFilter.
        :type ue_type: UeType
        """

        self._ue_type = ue_type

    @property
    def aerial_srv_dnn_ind(self):
        """Gets the aerial_srv_dnn_ind of this UeInAreaFilter.


        :return: The aerial_srv_dnn_ind of this UeInAreaFilter.
        :rtype: bool
        """
        return self._aerial_srv_dnn_ind

    @aerial_srv_dnn_ind.setter
    def aerial_srv_dnn_ind(self, aerial_srv_dnn_ind):
        """Sets the aerial_srv_dnn_ind of this UeInAreaFilter.


        :param aerial_srv_dnn_ind: The aerial_srv_dnn_ind of this UeInAreaFilter.
        :type aerial_srv_dnn_ind: bool
        """

        self._aerial_srv_dnn_ind = aerial_srv_dnn_ind
