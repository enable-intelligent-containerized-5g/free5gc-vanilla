# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.amf_event_trigger import AmfEventTrigger
from openapi_server.models.notification_flag import NotificationFlag
from openapi_server.models.partitioning_criteria import PartitioningCriteria
from openapi_server import util

from openapi_server.models.amf_event_trigger import AmfEventTrigger  # noqa: E501
from openapi_server.models.notification_flag import NotificationFlag  # noqa: E501
from openapi_server.models.partitioning_criteria import PartitioningCriteria  # noqa: E501

class AmfEventMode(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, trigger=None, max_reports=None, expiry=None, rep_period=None, samp_ratio=None, partitioning_criteria=None, notif_flag=None):  # noqa: E501
        """AmfEventMode - a model defined in OpenAPI

        :param trigger: The trigger of this AmfEventMode.  # noqa: E501
        :type trigger: AmfEventTrigger
        :param max_reports: The max_reports of this AmfEventMode.  # noqa: E501
        :type max_reports: int
        :param expiry: The expiry of this AmfEventMode.  # noqa: E501
        :type expiry: datetime
        :param rep_period: The rep_period of this AmfEventMode.  # noqa: E501
        :type rep_period: int
        :param samp_ratio: The samp_ratio of this AmfEventMode.  # noqa: E501
        :type samp_ratio: int
        :param partitioning_criteria: The partitioning_criteria of this AmfEventMode.  # noqa: E501
        :type partitioning_criteria: List[PartitioningCriteria]
        :param notif_flag: The notif_flag of this AmfEventMode.  # noqa: E501
        :type notif_flag: NotificationFlag
        """
        self.openapi_types = {
            'trigger': AmfEventTrigger,
            'max_reports': int,
            'expiry': datetime,
            'rep_period': int,
            'samp_ratio': int,
            'partitioning_criteria': List[PartitioningCriteria],
            'notif_flag': NotificationFlag
        }

        self.attribute_map = {
            'trigger': 'trigger',
            'max_reports': 'maxReports',
            'expiry': 'expiry',
            'rep_period': 'repPeriod',
            'samp_ratio': 'sampRatio',
            'partitioning_criteria': 'partitioningCriteria',
            'notif_flag': 'notifFlag'
        }

        self.trigger = trigger
        self.max_reports = max_reports
        self.expiry = expiry
        self.rep_period = rep_period
        self.samp_ratio = samp_ratio
        self.partitioning_criteria = partitioning_criteria
        self.notif_flag = notif_flag

    @classmethod
    def from_dict(cls, dikt) -> 'AmfEventMode':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The AmfEventMode of this AmfEventMode.  # noqa: E501
        :rtype: AmfEventMode
        """
        return util.deserialize_model(dikt, cls)

    @property
    def trigger(self):
        """Gets the trigger of this AmfEventMode.


        :return: The trigger of this AmfEventMode.
        :rtype: AmfEventTrigger
        """
        return self._trigger

    @trigger.setter
    def trigger(self, trigger):
        """Sets the trigger of this AmfEventMode.


        :param trigger: The trigger of this AmfEventMode.
        :type trigger: AmfEventTrigger
        """
        if trigger is None:
            raise ValueError("Invalid value for `trigger`, must not be `None`")  # noqa: E501

        self._trigger = trigger

    @property
    def max_reports(self):
        """Gets the max_reports of this AmfEventMode.


        :return: The max_reports of this AmfEventMode.
        :rtype: int
        """
        return self._max_reports

    @max_reports.setter
    def max_reports(self, max_reports):
        """Sets the max_reports of this AmfEventMode.


        :param max_reports: The max_reports of this AmfEventMode.
        :type max_reports: int
        """

        self._max_reports = max_reports

    @property
    def expiry(self):
        """Gets the expiry of this AmfEventMode.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :return: The expiry of this AmfEventMode.
        :rtype: datetime
        """
        return self._expiry

    @expiry.setter
    def expiry(self, expiry):
        """Sets the expiry of this AmfEventMode.

        string with format 'date-time' as defined in OpenAPI.  # noqa: E501

        :param expiry: The expiry of this AmfEventMode.
        :type expiry: datetime
        """

        self._expiry = expiry

    @property
    def rep_period(self):
        """Gets the rep_period of this AmfEventMode.

        indicating a time in seconds.  # noqa: E501

        :return: The rep_period of this AmfEventMode.
        :rtype: int
        """
        return self._rep_period

    @rep_period.setter
    def rep_period(self, rep_period):
        """Sets the rep_period of this AmfEventMode.

        indicating a time in seconds.  # noqa: E501

        :param rep_period: The rep_period of this AmfEventMode.
        :type rep_period: int
        """

        self._rep_period = rep_period

    @property
    def samp_ratio(self):
        """Gets the samp_ratio of this AmfEventMode.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :return: The samp_ratio of this AmfEventMode.
        :rtype: int
        """
        return self._samp_ratio

    @samp_ratio.setter
    def samp_ratio(self, samp_ratio):
        """Sets the samp_ratio of this AmfEventMode.

        Unsigned integer indicating Sampling Ratio (see clauses 4.15.1 of 3GPP TS 23.502), expressed in percent.    # noqa: E501

        :param samp_ratio: The samp_ratio of this AmfEventMode.
        :type samp_ratio: int
        """
        if samp_ratio is not None and samp_ratio > 100:  # noqa: E501
            raise ValueError("Invalid value for `samp_ratio`, must be a value less than or equal to `100`")  # noqa: E501
        if samp_ratio is not None and samp_ratio < 1:  # noqa: E501
            raise ValueError("Invalid value for `samp_ratio`, must be a value greater than or equal to `1`")  # noqa: E501

        self._samp_ratio = samp_ratio

    @property
    def partitioning_criteria(self):
        """Gets the partitioning_criteria of this AmfEventMode.


        :return: The partitioning_criteria of this AmfEventMode.
        :rtype: List[PartitioningCriteria]
        """
        return self._partitioning_criteria

    @partitioning_criteria.setter
    def partitioning_criteria(self, partitioning_criteria):
        """Sets the partitioning_criteria of this AmfEventMode.


        :param partitioning_criteria: The partitioning_criteria of this AmfEventMode.
        :type partitioning_criteria: List[PartitioningCriteria]
        """
        if partitioning_criteria is not None and len(partitioning_criteria) < 1:
            raise ValueError("Invalid value for `partitioning_criteria`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._partitioning_criteria = partitioning_criteria

    @property
    def notif_flag(self):
        """Gets the notif_flag of this AmfEventMode.


        :return: The notif_flag of this AmfEventMode.
        :rtype: NotificationFlag
        """
        return self._notif_flag

    @notif_flag.setter
    def notif_flag(self, notif_flag):
        """Sets the notif_flag of this AmfEventMode.


        :param notif_flag: The notif_flag of this AmfEventMode.
        :type notif_flag: NotificationFlag
        """

        self._notif_flag = notif_flag