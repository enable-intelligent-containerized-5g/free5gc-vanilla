# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class ScheduledCommunicationTime(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, days_of_week=None, time_of_day_start=None, time_of_day_end=None):  # noqa: E501
        """ScheduledCommunicationTime - a model defined in OpenAPI

        :param days_of_week: The days_of_week of this ScheduledCommunicationTime.  # noqa: E501
        :type days_of_week: List[int]
        :param time_of_day_start: The time_of_day_start of this ScheduledCommunicationTime.  # noqa: E501
        :type time_of_day_start: str
        :param time_of_day_end: The time_of_day_end of this ScheduledCommunicationTime.  # noqa: E501
        :type time_of_day_end: str
        """
        self.openapi_types = {
            'days_of_week': List[int],
            'time_of_day_start': str,
            'time_of_day_end': str
        }

        self.attribute_map = {
            'days_of_week': 'daysOfWeek',
            'time_of_day_start': 'timeOfDayStart',
            'time_of_day_end': 'timeOfDayEnd'
        }

        self.days_of_week = days_of_week
        self.time_of_day_start = time_of_day_start
        self.time_of_day_end = time_of_day_end

    @classmethod
    def from_dict(cls, dikt) -> 'ScheduledCommunicationTime':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ScheduledCommunicationTime of this ScheduledCommunicationTime.  # noqa: E501
        :rtype: ScheduledCommunicationTime
        """
        return util.deserialize_model(dikt, cls)

    @property
    def days_of_week(self):
        """Gets the days_of_week of this ScheduledCommunicationTime.

        Identifies the day(s) of the week. If absent, it indicates every day of the week.   # noqa: E501

        :return: The days_of_week of this ScheduledCommunicationTime.
        :rtype: List[int]
        """
        return self._days_of_week

    @days_of_week.setter
    def days_of_week(self, days_of_week):
        """Sets the days_of_week of this ScheduledCommunicationTime.

        Identifies the day(s) of the week. If absent, it indicates every day of the week.   # noqa: E501

        :param days_of_week: The days_of_week of this ScheduledCommunicationTime.
        :type days_of_week: List[int]
        """
        if days_of_week is not None and len(days_of_week) > 6:
            raise ValueError("Invalid value for `days_of_week`, number of items must be less than or equal to `6`")  # noqa: E501
        if days_of_week is not None and len(days_of_week) < 1:
            raise ValueError("Invalid value for `days_of_week`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._days_of_week = days_of_week

    @property
    def time_of_day_start(self):
        """Gets the time_of_day_start of this ScheduledCommunicationTime.

        String with format partial-time or full-time as defined in clause 5.6 of IETF RFC 3339. Examples, 20:15:00, 20:15:00-08:00 (for 8 hours behind UTC).    # noqa: E501

        :return: The time_of_day_start of this ScheduledCommunicationTime.
        :rtype: str
        """
        return self._time_of_day_start

    @time_of_day_start.setter
    def time_of_day_start(self, time_of_day_start):
        """Sets the time_of_day_start of this ScheduledCommunicationTime.

        String with format partial-time or full-time as defined in clause 5.6 of IETF RFC 3339. Examples, 20:15:00, 20:15:00-08:00 (for 8 hours behind UTC).    # noqa: E501

        :param time_of_day_start: The time_of_day_start of this ScheduledCommunicationTime.
        :type time_of_day_start: str
        """

        self._time_of_day_start = time_of_day_start

    @property
    def time_of_day_end(self):
        """Gets the time_of_day_end of this ScheduledCommunicationTime.

        String with format partial-time or full-time as defined in clause 5.6 of IETF RFC 3339. Examples, 20:15:00, 20:15:00-08:00 (for 8 hours behind UTC).    # noqa: E501

        :return: The time_of_day_end of this ScheduledCommunicationTime.
        :rtype: str
        """
        return self._time_of_day_end

    @time_of_day_end.setter
    def time_of_day_end(self, time_of_day_end):
        """Sets the time_of_day_end of this ScheduledCommunicationTime.

        String with format partial-time or full-time as defined in clause 5.6 of IETF RFC 3339. Examples, 20:15:00, 20:15:00-08:00 (for 8 hours behind UTC).    # noqa: E501

        :param time_of_day_end: The time_of_day_end of this ScheduledCommunicationTime.
        :type time_of_day_end: str
        """

        self._time_of_day_end = time_of_day_end
