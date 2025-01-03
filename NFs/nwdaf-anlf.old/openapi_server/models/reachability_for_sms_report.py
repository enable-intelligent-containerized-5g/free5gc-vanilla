# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.access_type import AccessType
from openapi_server import util

from openapi_server.models.access_type import AccessType  # noqa: E501

class ReachabilityForSmsReport(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, smsf_access_type=None, max_availability_time=None):  # noqa: E501
        """ReachabilityForSmsReport - a model defined in OpenAPI

        :param smsf_access_type: The smsf_access_type of this ReachabilityForSmsReport.  # noqa: E501
        :type smsf_access_type: AccessType
        :param max_availability_time: The max_availability_time of this ReachabilityForSmsReport.  # noqa: E501
        :type max_availability_time: datetime
        """
        self.openapi_types = {
            'smsf_access_type': AccessType,
            'max_availability_time': datetime
        }

        self.attribute_map = {
            'smsf_access_type': 'smsfAccessType',
            'max_availability_time': 'maxAvailabilityTime'
        }

        self.smsf_access_type = smsf_access_type
        self.max_availability_time = max_availability_time

    @classmethod
    def from_dict(cls, dikt) -> 'ReachabilityForSmsReport':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ReachabilityForSmsReport of this ReachabilityForSmsReport.  # noqa: E501
        :rtype: ReachabilityForSmsReport
        """
        return util.deserialize_model(dikt, cls)

    @property
    def smsf_access_type(self):
        """Gets the smsf_access_type of this ReachabilityForSmsReport.


        :return: The smsf_access_type of this ReachabilityForSmsReport.
        :rtype: AccessType
        """
        return self._smsf_access_type

    @smsf_access_type.setter
    def smsf_access_type(self, smsf_access_type):
        """Sets the smsf_access_type of this ReachabilityForSmsReport.


        :param smsf_access_type: The smsf_access_type of this ReachabilityForSmsReport.
        :type smsf_access_type: AccessType
        """
        if smsf_access_type is None:
            raise ValueError("Invalid value for `smsf_access_type`, must not be `None`")  # noqa: E501

        self._smsf_access_type = smsf_access_type

    @property
    def max_availability_time(self):
        """Gets the max_availability_time of this ReachabilityForSmsReport.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :return: The max_availability_time of this ReachabilityForSmsReport.
        :rtype: datetime
        """
        return self._max_availability_time

    @max_availability_time.setter
    def max_availability_time(self, max_availability_time):
        """Sets the max_availability_time of this ReachabilityForSmsReport.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :param max_availability_time: The max_availability_time of this ReachabilityForSmsReport.
        :type max_availability_time: datetime
        """

        self._max_availability_time = max_availability_time
