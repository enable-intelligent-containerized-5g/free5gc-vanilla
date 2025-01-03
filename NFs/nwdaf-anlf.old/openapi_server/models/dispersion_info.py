# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.dispersion_collection import DispersionCollection
from openapi_server.models.dispersion_type import DispersionType
from openapi_server import util

from openapi_server.models.dispersion_collection import DispersionCollection  # noqa: E501
from openapi_server.models.dispersion_type import DispersionType  # noqa: E501

class DispersionInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ts_start=None, ts_duration=None, disper_collects=None, disper_type=None):  # noqa: E501
        """DispersionInfo - a model defined in OpenAPI

        :param ts_start: The ts_start of this DispersionInfo.  # noqa: E501
        :type ts_start: datetime
        :param ts_duration: The ts_duration of this DispersionInfo.  # noqa: E501
        :type ts_duration: int
        :param disper_collects: The disper_collects of this DispersionInfo.  # noqa: E501
        :type disper_collects: List[DispersionCollection]
        :param disper_type: The disper_type of this DispersionInfo.  # noqa: E501
        :type disper_type: DispersionType
        """
        self.openapi_types = {
            'ts_start': datetime,
            'ts_duration': int,
            'disper_collects': List[DispersionCollection],
            'disper_type': DispersionType
        }

        self.attribute_map = {
            'ts_start': 'tsStart',
            'ts_duration': 'tsDuration',
            'disper_collects': 'disperCollects',
            'disper_type': 'disperType'
        }

        self.ts_start = ts_start
        self.ts_duration = ts_duration
        self.disper_collects = disper_collects
        self.disper_type = disper_type

    @classmethod
    def from_dict(cls, dikt) -> 'DispersionInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The DispersionInfo of this DispersionInfo.  # noqa: E501
        :rtype: DispersionInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ts_start(self):
        """Gets the ts_start of this DispersionInfo.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :return: The ts_start of this DispersionInfo.
        :rtype: datetime
        """
        return self._ts_start

    @ts_start.setter
    def ts_start(self, ts_start):
        """Sets the ts_start of this DispersionInfo.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :param ts_start: The ts_start of this DispersionInfo.
        :type ts_start: datetime
        """
        if ts_start is None:
            raise ValueError("Invalid value for `ts_start`, must not be `None`")  # noqa: E501

        self._ts_start = ts_start

    @property
    def ts_duration(self):
        """Gets the ts_duration of this DispersionInfo.

        indicating a time in seconds.  # noqa: E501

        :return: The ts_duration of this DispersionInfo.
        :rtype: int
        """
        return self._ts_duration

    @ts_duration.setter
    def ts_duration(self, ts_duration):
        """Sets the ts_duration of this DispersionInfo.

        indicating a time in seconds.  # noqa: E501

        :param ts_duration: The ts_duration of this DispersionInfo.
        :type ts_duration: int
        """
        if ts_duration is None:
            raise ValueError("Invalid value for `ts_duration`, must not be `None`")  # noqa: E501

        self._ts_duration = ts_duration

    @property
    def disper_collects(self):
        """Gets the disper_collects of this DispersionInfo.


        :return: The disper_collects of this DispersionInfo.
        :rtype: List[DispersionCollection]
        """
        return self._disper_collects

    @disper_collects.setter
    def disper_collects(self, disper_collects):
        """Sets the disper_collects of this DispersionInfo.


        :param disper_collects: The disper_collects of this DispersionInfo.
        :type disper_collects: List[DispersionCollection]
        """
        if disper_collects is None:
            raise ValueError("Invalid value for `disper_collects`, must not be `None`")  # noqa: E501
        if disper_collects is not None and len(disper_collects) < 1:
            raise ValueError("Invalid value for `disper_collects`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._disper_collects = disper_collects

    @property
    def disper_type(self):
        """Gets the disper_type of this DispersionInfo.


        :return: The disper_type of this DispersionInfo.
        :rtype: DispersionType
        """
        return self._disper_type

    @disper_type.setter
    def disper_type(self, disper_type):
        """Sets the disper_type of this DispersionInfo.


        :param disper_type: The disper_type of this DispersionInfo.
        :type disper_type: DispersionType
        """
        if disper_type is None:
            raise ValueError("Invalid value for `disper_type`, must not be `None`")  # noqa: E501

        self._disper_type = disper_type
