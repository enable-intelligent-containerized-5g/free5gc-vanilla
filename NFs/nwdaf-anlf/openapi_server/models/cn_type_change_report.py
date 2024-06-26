# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.cn_type import CnType
from openapi_server import util

from openapi_server.models.cn_type import CnType  # noqa: E501

class CnTypeChangeReport(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, new_cn_type=None, old_cn_type=None):  # noqa: E501
        """CnTypeChangeReport - a model defined in OpenAPI

        :param new_cn_type: The new_cn_type of this CnTypeChangeReport.  # noqa: E501
        :type new_cn_type: CnType
        :param old_cn_type: The old_cn_type of this CnTypeChangeReport.  # noqa: E501
        :type old_cn_type: CnType
        """
        self.openapi_types = {
            'new_cn_type': CnType,
            'old_cn_type': CnType
        }

        self.attribute_map = {
            'new_cn_type': 'newCnType',
            'old_cn_type': 'oldCnType'
        }

        self.new_cn_type = new_cn_type
        self.old_cn_type = old_cn_type

    @classmethod
    def from_dict(cls, dikt) -> 'CnTypeChangeReport':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The CnTypeChangeReport of this CnTypeChangeReport.  # noqa: E501
        :rtype: CnTypeChangeReport
        """
        return util.deserialize_model(dikt, cls)

    @property
    def new_cn_type(self):
        """Gets the new_cn_type of this CnTypeChangeReport.


        :return: The new_cn_type of this CnTypeChangeReport.
        :rtype: CnType
        """
        return self._new_cn_type

    @new_cn_type.setter
    def new_cn_type(self, new_cn_type):
        """Sets the new_cn_type of this CnTypeChangeReport.


        :param new_cn_type: The new_cn_type of this CnTypeChangeReport.
        :type new_cn_type: CnType
        """
        if new_cn_type is None:
            raise ValueError("Invalid value for `new_cn_type`, must not be `None`")  # noqa: E501

        self._new_cn_type = new_cn_type

    @property
    def old_cn_type(self):
        """Gets the old_cn_type of this CnTypeChangeReport.


        :return: The old_cn_type of this CnTypeChangeReport.
        :rtype: CnType
        """
        return self._old_cn_type

    @old_cn_type.setter
    def old_cn_type(self, old_cn_type):
        """Sets the old_cn_type of this CnTypeChangeReport.


        :param old_cn_type: The old_cn_type of this CnTypeChangeReport.
        :type old_cn_type: CnType
        """

        self._old_cn_type = old_cn_type