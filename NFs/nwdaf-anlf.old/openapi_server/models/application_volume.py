# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class ApplicationVolume(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, app_id=None, app_volume=None):  # noqa: E501
        """ApplicationVolume - a model defined in OpenAPI

        :param app_id: The app_id of this ApplicationVolume.  # noqa: E501
        :type app_id: str
        :param app_volume: The app_volume of this ApplicationVolume.  # noqa: E501
        :type app_volume: int
        """
        self.openapi_types = {
            'app_id': str,
            'app_volume': int
        }

        self.attribute_map = {
            'app_id': 'appId',
            'app_volume': 'appVolume'
        }

        self.app_id = app_id
        self.app_volume = app_volume

    @classmethod
    def from_dict(cls, dikt) -> 'ApplicationVolume':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The ApplicationVolume of this ApplicationVolume.  # noqa: E501
        :rtype: ApplicationVolume
        """
        return util.deserialize_model(dikt, cls)

    @property
    def app_id(self):
        """Gets the app_id of this ApplicationVolume.

        String providing an application identifier.  # noqa: E501

        :return: The app_id of this ApplicationVolume.
        :rtype: str
        """
        return self._app_id

    @app_id.setter
    def app_id(self, app_id):
        """Sets the app_id of this ApplicationVolume.

        String providing an application identifier.  # noqa: E501

        :param app_id: The app_id of this ApplicationVolume.
        :type app_id: str
        """
        if app_id is None:
            raise ValueError("Invalid value for `app_id`, must not be `None`")  # noqa: E501

        self._app_id = app_id

    @property
    def app_volume(self):
        """Gets the app_volume of this ApplicationVolume.

        Unsigned integer identifying a volume in units of bytes.  # noqa: E501

        :return: The app_volume of this ApplicationVolume.
        :rtype: int
        """
        return self._app_volume

    @app_volume.setter
    def app_volume(self, app_volume):
        """Sets the app_volume of this ApplicationVolume.

        Unsigned integer identifying a volume in units of bytes.  # noqa: E501

        :param app_volume: The app_volume of this ApplicationVolume.
        :type app_volume: int
        """
        if app_volume is None:
            raise ValueError("Invalid value for `app_volume`, must not be `None`")  # noqa: E501
        if app_volume is not None and app_volume < 0:  # noqa: E501
            raise ValueError("Invalid value for `app_volume`, must be a value greater than or equal to `0`")  # noqa: E501

        self._app_volume = app_volume
