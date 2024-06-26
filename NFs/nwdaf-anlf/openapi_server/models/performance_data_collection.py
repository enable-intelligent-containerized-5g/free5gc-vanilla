# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.addr_fqdn import AddrFqdn
from openapi_server.models.flow_info import FlowInfo
from openapi_server.models.ip_addr import IpAddr
from openapi_server.models.location_area5_g import LocationArea5G
from openapi_server.models.performance_data import PerformanceData
from openapi_server import util

from openapi_server.models.addr_fqdn import AddrFqdn  # noqa: E501
from openapi_server.models.flow_info import FlowInfo  # noqa: E501
from openapi_server.models.ip_addr import IpAddr  # noqa: E501
from openapi_server.models.location_area5_g import LocationArea5G  # noqa: E501
from openapi_server.models.performance_data import PerformanceData  # noqa: E501

class PerformanceDataCollection(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, app_id=None, ue_ip_addr=None, ip_traffic_filter=None, ue_loc=None, app_locs=None, as_addr=None, perf_data=None, time_stamp=None):  # noqa: E501
        """PerformanceDataCollection - a model defined in OpenAPI

        :param app_id: The app_id of this PerformanceDataCollection.  # noqa: E501
        :type app_id: str
        :param ue_ip_addr: The ue_ip_addr of this PerformanceDataCollection.  # noqa: E501
        :type ue_ip_addr: IpAddr
        :param ip_traffic_filter: The ip_traffic_filter of this PerformanceDataCollection.  # noqa: E501
        :type ip_traffic_filter: FlowInfo
        :param ue_loc: The ue_loc of this PerformanceDataCollection.  # noqa: E501
        :type ue_loc: LocationArea5G
        :param app_locs: The app_locs of this PerformanceDataCollection.  # noqa: E501
        :type app_locs: List[str]
        :param as_addr: The as_addr of this PerformanceDataCollection.  # noqa: E501
        :type as_addr: AddrFqdn
        :param perf_data: The perf_data of this PerformanceDataCollection.  # noqa: E501
        :type perf_data: PerformanceData
        :param time_stamp: The time_stamp of this PerformanceDataCollection.  # noqa: E501
        :type time_stamp: datetime
        """
        self.openapi_types = {
            'app_id': str,
            'ue_ip_addr': IpAddr,
            'ip_traffic_filter': FlowInfo,
            'ue_loc': LocationArea5G,
            'app_locs': List[str],
            'as_addr': AddrFqdn,
            'perf_data': PerformanceData,
            'time_stamp': datetime
        }

        self.attribute_map = {
            'app_id': 'appId',
            'ue_ip_addr': 'ueIpAddr',
            'ip_traffic_filter': 'ipTrafficFilter',
            'ue_loc': 'ueLoc',
            'app_locs': 'appLocs',
            'as_addr': 'asAddr',
            'perf_data': 'perfData',
            'time_stamp': 'timeStamp'
        }

        self.app_id = app_id
        self.ue_ip_addr = ue_ip_addr
        self.ip_traffic_filter = ip_traffic_filter
        self.ue_loc = ue_loc
        self.app_locs = app_locs
        self.as_addr = as_addr
        self.perf_data = perf_data
        self.time_stamp = time_stamp

    @classmethod
    def from_dict(cls, dikt) -> 'PerformanceDataCollection':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PerformanceDataCollection of this PerformanceDataCollection.  # noqa: E501
        :rtype: PerformanceDataCollection
        """
        return util.deserialize_model(dikt, cls)

    @property
    def app_id(self):
        """Gets the app_id of this PerformanceDataCollection.

        String providing an application identifier.  # noqa: E501

        :return: The app_id of this PerformanceDataCollection.
        :rtype: str
        """
        return self._app_id

    @app_id.setter
    def app_id(self, app_id):
        """Sets the app_id of this PerformanceDataCollection.

        String providing an application identifier.  # noqa: E501

        :param app_id: The app_id of this PerformanceDataCollection.
        :type app_id: str
        """

        self._app_id = app_id

    @property
    def ue_ip_addr(self):
        """Gets the ue_ip_addr of this PerformanceDataCollection.


        :return: The ue_ip_addr of this PerformanceDataCollection.
        :rtype: IpAddr
        """
        return self._ue_ip_addr

    @ue_ip_addr.setter
    def ue_ip_addr(self, ue_ip_addr):
        """Sets the ue_ip_addr of this PerformanceDataCollection.


        :param ue_ip_addr: The ue_ip_addr of this PerformanceDataCollection.
        :type ue_ip_addr: IpAddr
        """

        self._ue_ip_addr = ue_ip_addr

    @property
    def ip_traffic_filter(self):
        """Gets the ip_traffic_filter of this PerformanceDataCollection.


        :return: The ip_traffic_filter of this PerformanceDataCollection.
        :rtype: FlowInfo
        """
        return self._ip_traffic_filter

    @ip_traffic_filter.setter
    def ip_traffic_filter(self, ip_traffic_filter):
        """Sets the ip_traffic_filter of this PerformanceDataCollection.


        :param ip_traffic_filter: The ip_traffic_filter of this PerformanceDataCollection.
        :type ip_traffic_filter: FlowInfo
        """

        self._ip_traffic_filter = ip_traffic_filter

    @property
    def ue_loc(self):
        """Gets the ue_loc of this PerformanceDataCollection.


        :return: The ue_loc of this PerformanceDataCollection.
        :rtype: LocationArea5G
        """
        return self._ue_loc

    @ue_loc.setter
    def ue_loc(self, ue_loc):
        """Sets the ue_loc of this PerformanceDataCollection.


        :param ue_loc: The ue_loc of this PerformanceDataCollection.
        :type ue_loc: LocationArea5G
        """

        self._ue_loc = ue_loc

    @property
    def app_locs(self):
        """Gets the app_locs of this PerformanceDataCollection.


        :return: The app_locs of this PerformanceDataCollection.
        :rtype: List[str]
        """
        return self._app_locs

    @app_locs.setter
    def app_locs(self, app_locs):
        """Sets the app_locs of this PerformanceDataCollection.


        :param app_locs: The app_locs of this PerformanceDataCollection.
        :type app_locs: List[str]
        """
        if app_locs is not None and len(app_locs) < 1:
            raise ValueError("Invalid value for `app_locs`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._app_locs = app_locs

    @property
    def as_addr(self):
        """Gets the as_addr of this PerformanceDataCollection.


        :return: The as_addr of this PerformanceDataCollection.
        :rtype: AddrFqdn
        """
        return self._as_addr

    @as_addr.setter
    def as_addr(self, as_addr):
        """Sets the as_addr of this PerformanceDataCollection.


        :param as_addr: The as_addr of this PerformanceDataCollection.
        :type as_addr: AddrFqdn
        """

        self._as_addr = as_addr

    @property
    def perf_data(self):
        """Gets the perf_data of this PerformanceDataCollection.


        :return: The perf_data of this PerformanceDataCollection.
        :rtype: PerformanceData
        """
        return self._perf_data

    @perf_data.setter
    def perf_data(self, perf_data):
        """Sets the perf_data of this PerformanceDataCollection.


        :param perf_data: The perf_data of this PerformanceDataCollection.
        :type perf_data: PerformanceData
        """
        if perf_data is None:
            raise ValueError("Invalid value for `perf_data`, must not be `None`")  # noqa: E501

        self._perf_data = perf_data

    @property
    def time_stamp(self):
        """Gets the time_stamp of this PerformanceDataCollection.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :return: The time_stamp of this PerformanceDataCollection.
        :rtype: datetime
        """
        return self._time_stamp

    @time_stamp.setter
    def time_stamp(self, time_stamp):
        """Sets the time_stamp of this PerformanceDataCollection.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :param time_stamp: The time_stamp of this PerformanceDataCollection.
        :type time_stamp: datetime
        """
        if time_stamp is None:
            raise ValueError("Invalid value for `time_stamp`, must not be `None`")  # noqa: E501

        self._time_stamp = time_stamp