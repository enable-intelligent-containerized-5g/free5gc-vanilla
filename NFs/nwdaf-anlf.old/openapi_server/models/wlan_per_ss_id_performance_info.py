# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.wlan_per_ts_performance_info import WlanPerTsPerformanceInfo
from openapi_server import util

from openapi_server.models.wlan_per_ts_performance_info import WlanPerTsPerformanceInfo  # noqa: E501

class WlanPerSsIdPerformanceInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ss_id=None, wlan_per_ts_infos=None):  # noqa: E501
        """WlanPerSsIdPerformanceInfo - a model defined in OpenAPI

        :param ss_id: The ss_id of this WlanPerSsIdPerformanceInfo.  # noqa: E501
        :type ss_id: str
        :param wlan_per_ts_infos: The wlan_per_ts_infos of this WlanPerSsIdPerformanceInfo.  # noqa: E501
        :type wlan_per_ts_infos: List[WlanPerTsPerformanceInfo]
        """
        self.openapi_types = {
            'ss_id': str,
            'wlan_per_ts_infos': List[WlanPerTsPerformanceInfo]
        }

        self.attribute_map = {
            'ss_id': 'ssId',
            'wlan_per_ts_infos': 'wlanPerTsInfos'
        }

        self.ss_id = ss_id
        self.wlan_per_ts_infos = wlan_per_ts_infos

    @classmethod
    def from_dict(cls, dikt) -> 'WlanPerSsIdPerformanceInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The WlanPerSsIdPerformanceInfo of this WlanPerSsIdPerformanceInfo.  # noqa: E501
        :rtype: WlanPerSsIdPerformanceInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ss_id(self):
        """Gets the ss_id of this WlanPerSsIdPerformanceInfo.


        :return: The ss_id of this WlanPerSsIdPerformanceInfo.
        :rtype: str
        """
        return self._ss_id

    @ss_id.setter
    def ss_id(self, ss_id):
        """Sets the ss_id of this WlanPerSsIdPerformanceInfo.


        :param ss_id: The ss_id of this WlanPerSsIdPerformanceInfo.
        :type ss_id: str
        """
        if ss_id is None:
            raise ValueError("Invalid value for `ss_id`, must not be `None`")  # noqa: E501

        self._ss_id = ss_id

    @property
    def wlan_per_ts_infos(self):
        """Gets the wlan_per_ts_infos of this WlanPerSsIdPerformanceInfo.


        :return: The wlan_per_ts_infos of this WlanPerSsIdPerformanceInfo.
        :rtype: List[WlanPerTsPerformanceInfo]
        """
        return self._wlan_per_ts_infos

    @wlan_per_ts_infos.setter
    def wlan_per_ts_infos(self, wlan_per_ts_infos):
        """Sets the wlan_per_ts_infos of this WlanPerSsIdPerformanceInfo.


        :param wlan_per_ts_infos: The wlan_per_ts_infos of this WlanPerSsIdPerformanceInfo.
        :type wlan_per_ts_infos: List[WlanPerTsPerformanceInfo]
        """
        if wlan_per_ts_infos is None:
            raise ValueError("Invalid value for `wlan_per_ts_infos`, must not be `None`")  # noqa: E501
        if wlan_per_ts_infos is not None and len(wlan_per_ts_infos) < 1:
            raise ValueError("Invalid value for `wlan_per_ts_infos`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._wlan_per_ts_infos = wlan_per_ts_infos
