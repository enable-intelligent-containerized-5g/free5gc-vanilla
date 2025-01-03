# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.smcce_ue_list import SmcceUeList
from openapi_server.models.snssai import Snssai
from openapi_server import util

from openapi_server.models.smcce_ue_list import SmcceUeList  # noqa: E501
from openapi_server.models.snssai import Snssai  # noqa: E501

class SmcceInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, dnn=None, snssai=None, smcce_ue_list=None):  # noqa: E501
        """SmcceInfo - a model defined in OpenAPI

        :param dnn: The dnn of this SmcceInfo.  # noqa: E501
        :type dnn: str
        :param snssai: The snssai of this SmcceInfo.  # noqa: E501
        :type snssai: Snssai
        :param smcce_ue_list: The smcce_ue_list of this SmcceInfo.  # noqa: E501
        :type smcce_ue_list: SmcceUeList
        """
        self.openapi_types = {
            'dnn': str,
            'snssai': Snssai,
            'smcce_ue_list': SmcceUeList
        }

        self.attribute_map = {
            'dnn': 'dnn',
            'snssai': 'snssai',
            'smcce_ue_list': 'smcceUeList'
        }

        self.dnn = dnn
        self.snssai = snssai
        self.smcce_ue_list = smcce_ue_list

    @classmethod
    def from_dict(cls, dikt) -> 'SmcceInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The SmcceInfo of this SmcceInfo.  # noqa: E501
        :rtype: SmcceInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def dnn(self):
        """Gets the dnn of this SmcceInfo.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :return: The dnn of this SmcceInfo.
        :rtype: str
        """
        return self._dnn

    @dnn.setter
    def dnn(self, dnn):
        """Sets the dnn of this SmcceInfo.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :param dnn: The dnn of this SmcceInfo.
        :type dnn: str
        """

        self._dnn = dnn

    @property
    def snssai(self):
        """Gets the snssai of this SmcceInfo.


        :return: The snssai of this SmcceInfo.
        :rtype: Snssai
        """
        return self._snssai

    @snssai.setter
    def snssai(self, snssai):
        """Sets the snssai of this SmcceInfo.


        :param snssai: The snssai of this SmcceInfo.
        :type snssai: Snssai
        """

        self._snssai = snssai

    @property
    def smcce_ue_list(self):
        """Gets the smcce_ue_list of this SmcceInfo.


        :return: The smcce_ue_list of this SmcceInfo.
        :rtype: SmcceUeList
        """
        return self._smcce_ue_list

    @smcce_ue_list.setter
    def smcce_ue_list(self, smcce_ue_list):
        """Sets the smcce_ue_list of this SmcceInfo.


        :param smcce_ue_list: The smcce_ue_list of this SmcceInfo.
        :type smcce_ue_list: SmcceUeList
        """
        if smcce_ue_list is None:
            raise ValueError("Invalid value for `smcce_ue_list`, must not be `None`")  # noqa: E501

        self._smcce_ue_list = smcce_ue_list
