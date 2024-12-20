# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.per_ue_attribute import PerUeAttribute
from openapi_server import util

from openapi_server.models.per_ue_attribute import PerUeAttribute  # noqa: E501

class CollectiveBehaviourInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, col_attrib=None, no_of_ues=None, app_ids=None, ext_ue_ids=None, ue_ids=None):  # noqa: E501
        """CollectiveBehaviourInfo - a model defined in OpenAPI

        :param col_attrib: The col_attrib of this CollectiveBehaviourInfo.  # noqa: E501
        :type col_attrib: List[PerUeAttribute]
        :param no_of_ues: The no_of_ues of this CollectiveBehaviourInfo.  # noqa: E501
        :type no_of_ues: int
        :param app_ids: The app_ids of this CollectiveBehaviourInfo.  # noqa: E501
        :type app_ids: List[str]
        :param ext_ue_ids: The ext_ue_ids of this CollectiveBehaviourInfo.  # noqa: E501
        :type ext_ue_ids: List[str]
        :param ue_ids: The ue_ids of this CollectiveBehaviourInfo.  # noqa: E501
        :type ue_ids: List[str]
        """
        self.openapi_types = {
            'col_attrib': List[PerUeAttribute],
            'no_of_ues': int,
            'app_ids': List[str],
            'ext_ue_ids': List[str],
            'ue_ids': List[str]
        }

        self.attribute_map = {
            'col_attrib': 'colAttrib',
            'no_of_ues': 'noOfUes',
            'app_ids': 'appIds',
            'ext_ue_ids': 'extUeIds',
            'ue_ids': 'ueIds'
        }

        self.col_attrib = col_attrib
        self.no_of_ues = no_of_ues
        self.app_ids = app_ids
        self.ext_ue_ids = ext_ue_ids
        self.ue_ids = ue_ids

    @classmethod
    def from_dict(cls, dikt) -> 'CollectiveBehaviourInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The CollectiveBehaviourInfo of this CollectiveBehaviourInfo.  # noqa: E501
        :rtype: CollectiveBehaviourInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def col_attrib(self):
        """Gets the col_attrib of this CollectiveBehaviourInfo.


        :return: The col_attrib of this CollectiveBehaviourInfo.
        :rtype: List[PerUeAttribute]
        """
        return self._col_attrib

    @col_attrib.setter
    def col_attrib(self, col_attrib):
        """Sets the col_attrib of this CollectiveBehaviourInfo.


        :param col_attrib: The col_attrib of this CollectiveBehaviourInfo.
        :type col_attrib: List[PerUeAttribute]
        """
        if col_attrib is None:
            raise ValueError("Invalid value for `col_attrib`, must not be `None`")  # noqa: E501
        if col_attrib is not None and len(col_attrib) < 1:
            raise ValueError("Invalid value for `col_attrib`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._col_attrib = col_attrib

    @property
    def no_of_ues(self):
        """Gets the no_of_ues of this CollectiveBehaviourInfo.

        Total number of UEs that fulfil a collective within the area of interest.  # noqa: E501

        :return: The no_of_ues of this CollectiveBehaviourInfo.
        :rtype: int
        """
        return self._no_of_ues

    @no_of_ues.setter
    def no_of_ues(self, no_of_ues):
        """Sets the no_of_ues of this CollectiveBehaviourInfo.

        Total number of UEs that fulfil a collective within the area of interest.  # noqa: E501

        :param no_of_ues: The no_of_ues of this CollectiveBehaviourInfo.
        :type no_of_ues: int
        """

        self._no_of_ues = no_of_ues

    @property
    def app_ids(self):
        """Gets the app_ids of this CollectiveBehaviourInfo.


        :return: The app_ids of this CollectiveBehaviourInfo.
        :rtype: List[str]
        """
        return self._app_ids

    @app_ids.setter
    def app_ids(self, app_ids):
        """Sets the app_ids of this CollectiveBehaviourInfo.


        :param app_ids: The app_ids of this CollectiveBehaviourInfo.
        :type app_ids: List[str]
        """
        if app_ids is not None and len(app_ids) < 1:
            raise ValueError("Invalid value for `app_ids`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._app_ids = app_ids

    @property
    def ext_ue_ids(self):
        """Gets the ext_ue_ids of this CollectiveBehaviourInfo.


        :return: The ext_ue_ids of this CollectiveBehaviourInfo.
        :rtype: List[str]
        """
        return self._ext_ue_ids

    @ext_ue_ids.setter
    def ext_ue_ids(self, ext_ue_ids):
        """Sets the ext_ue_ids of this CollectiveBehaviourInfo.


        :param ext_ue_ids: The ext_ue_ids of this CollectiveBehaviourInfo.
        :type ext_ue_ids: List[str]
        """
        if ext_ue_ids is not None and len(ext_ue_ids) < 1:
            raise ValueError("Invalid value for `ext_ue_ids`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ext_ue_ids = ext_ue_ids

    @property
    def ue_ids(self):
        """Gets the ue_ids of this CollectiveBehaviourInfo.


        :return: The ue_ids of this CollectiveBehaviourInfo.
        :rtype: List[str]
        """
        return self._ue_ids

    @ue_ids.setter
    def ue_ids(self, ue_ids):
        """Sets the ue_ids of this CollectiveBehaviourInfo.


        :param ue_ids: The ue_ids of this CollectiveBehaviourInfo.
        :type ue_ids: List[str]
        """
        if ue_ids is not None and len(ue_ids) < 1:
            raise ValueError("Invalid value for `ue_ids`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ue_ids = ue_ids
