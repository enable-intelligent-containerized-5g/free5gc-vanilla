# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class LossConnectivityCfg(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, max_detection_time=None):  # noqa: E501
        """LossConnectivityCfg - a model defined in OpenAPI

        :param max_detection_time: The max_detection_time of this LossConnectivityCfg.  # noqa: E501
        :type max_detection_time: int
        """
        self.openapi_types = {
            'max_detection_time': int
        }

        self.attribute_map = {
            'max_detection_time': 'maxDetectionTime'
        }

        self.max_detection_time = max_detection_time

    @classmethod
    def from_dict(cls, dikt) -> 'LossConnectivityCfg':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The LossConnectivityCfg of this LossConnectivityCfg.  # noqa: E501
        :rtype: LossConnectivityCfg
        """
        return util.deserialize_model(dikt, cls)

    @property
    def max_detection_time(self):
        """Gets the max_detection_time of this LossConnectivityCfg.

        indicating a time in seconds.  # noqa: E501

        :return: The max_detection_time of this LossConnectivityCfg.
        :rtype: int
        """
        return self._max_detection_time

    @max_detection_time.setter
    def max_detection_time(self, max_detection_time):
        """Sets the max_detection_time of this LossConnectivityCfg.

        indicating a time in seconds.  # noqa: E501

        :param max_detection_time: The max_detection_time of this LossConnectivityCfg.
        :type max_detection_time: int
        """

        self._max_detection_time = max_detection_time
