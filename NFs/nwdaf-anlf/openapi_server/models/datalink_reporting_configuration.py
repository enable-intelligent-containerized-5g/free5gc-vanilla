# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ddd_traffic_descriptor import DddTrafficDescriptor
from openapi_server.models.dl_data_delivery_status import DlDataDeliveryStatus
from openapi_server.models.snssai import Snssai
from openapi_server import util

from openapi_server.models.ddd_traffic_descriptor import DddTrafficDescriptor  # noqa: E501
from openapi_server.models.dl_data_delivery_status import DlDataDeliveryStatus  # noqa: E501
from openapi_server.models.snssai import Snssai  # noqa: E501

class DatalinkReportingConfiguration(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ddd_traffic_des=None, dnn=None, slice=None, ddd_status_list=None):  # noqa: E501
        """DatalinkReportingConfiguration - a model defined in OpenAPI

        :param ddd_traffic_des: The ddd_traffic_des of this DatalinkReportingConfiguration.  # noqa: E501
        :type ddd_traffic_des: List[DddTrafficDescriptor]
        :param dnn: The dnn of this DatalinkReportingConfiguration.  # noqa: E501
        :type dnn: str
        :param slice: The slice of this DatalinkReportingConfiguration.  # noqa: E501
        :type slice: Snssai
        :param ddd_status_list: The ddd_status_list of this DatalinkReportingConfiguration.  # noqa: E501
        :type ddd_status_list: List[DlDataDeliveryStatus]
        """
        self.openapi_types = {
            'ddd_traffic_des': List[DddTrafficDescriptor],
            'dnn': str,
            'slice': Snssai,
            'ddd_status_list': List[DlDataDeliveryStatus]
        }

        self.attribute_map = {
            'ddd_traffic_des': 'dddTrafficDes',
            'dnn': 'dnn',
            'slice': 'slice',
            'ddd_status_list': 'dddStatusList'
        }

        self.ddd_traffic_des = ddd_traffic_des
        self.dnn = dnn
        self.slice = slice
        self.ddd_status_list = ddd_status_list

    @classmethod
    def from_dict(cls, dikt) -> 'DatalinkReportingConfiguration':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The DatalinkReportingConfiguration of this DatalinkReportingConfiguration.  # noqa: E501
        :rtype: DatalinkReportingConfiguration
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ddd_traffic_des(self):
        """Gets the ddd_traffic_des of this DatalinkReportingConfiguration.


        :return: The ddd_traffic_des of this DatalinkReportingConfiguration.
        :rtype: List[DddTrafficDescriptor]
        """
        return self._ddd_traffic_des

    @ddd_traffic_des.setter
    def ddd_traffic_des(self, ddd_traffic_des):
        """Sets the ddd_traffic_des of this DatalinkReportingConfiguration.


        :param ddd_traffic_des: The ddd_traffic_des of this DatalinkReportingConfiguration.
        :type ddd_traffic_des: List[DddTrafficDescriptor]
        """
        if ddd_traffic_des is not None and len(ddd_traffic_des) < 1:
            raise ValueError("Invalid value for `ddd_traffic_des`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ddd_traffic_des = ddd_traffic_des

    @property
    def dnn(self):
        """Gets the dnn of this DatalinkReportingConfiguration.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :return: The dnn of this DatalinkReportingConfiguration.
        :rtype: str
        """
        return self._dnn

    @dnn.setter
    def dnn(self, dnn):
        """Sets the dnn of this DatalinkReportingConfiguration.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :param dnn: The dnn of this DatalinkReportingConfiguration.
        :type dnn: str
        """

        self._dnn = dnn

    @property
    def slice(self):
        """Gets the slice of this DatalinkReportingConfiguration.


        :return: The slice of this DatalinkReportingConfiguration.
        :rtype: Snssai
        """
        return self._slice

    @slice.setter
    def slice(self, slice):
        """Sets the slice of this DatalinkReportingConfiguration.


        :param slice: The slice of this DatalinkReportingConfiguration.
        :type slice: Snssai
        """

        self._slice = slice

    @property
    def ddd_status_list(self):
        """Gets the ddd_status_list of this DatalinkReportingConfiguration.


        :return: The ddd_status_list of this DatalinkReportingConfiguration.
        :rtype: List[DlDataDeliveryStatus]
        """
        return self._ddd_status_list

    @ddd_status_list.setter
    def ddd_status_list(self, ddd_status_list):
        """Sets the ddd_status_list of this DatalinkReportingConfiguration.


        :param ddd_status_list: The ddd_status_list of this DatalinkReportingConfiguration.
        :type ddd_status_list: List[DlDataDeliveryStatus]
        """
        if ddd_status_list is not None and len(ddd_status_list) < 1:
            raise ValueError("Invalid value for `ddd_status_list`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ddd_status_list = ddd_status_list