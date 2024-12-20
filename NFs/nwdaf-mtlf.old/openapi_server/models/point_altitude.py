# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.gad_shape import GADShape
from openapi_server.models.geographical_coordinates import GeographicalCoordinates
from openapi_server.models.supported_gad_shapes import SupportedGADShapes
from openapi_server import util

from openapi_server.models.gad_shape import GADShape  # noqa: E501
from openapi_server.models.geographical_coordinates import GeographicalCoordinates  # noqa: E501
from openapi_server.models.supported_gad_shapes import SupportedGADShapes  # noqa: E501

class PointAltitude(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, shape=None, point=None, altitude=None):  # noqa: E501
        """PointAltitude - a model defined in OpenAPI

        :param shape: The shape of this PointAltitude.  # noqa: E501
        :type shape: SupportedGADShapes
        :param point: The point of this PointAltitude.  # noqa: E501
        :type point: GeographicalCoordinates
        :param altitude: The altitude of this PointAltitude.  # noqa: E501
        :type altitude: float
        """
        self.openapi_types = {
            'shape': SupportedGADShapes,
            'point': GeographicalCoordinates,
            'altitude': float
        }

        self.attribute_map = {
            'shape': 'shape',
            'point': 'point',
            'altitude': 'altitude'
        }

        self.shape = shape
        self.point = point
        self.altitude = altitude

    @classmethod
    def from_dict(cls, dikt) -> 'PointAltitude':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PointAltitude of this PointAltitude.  # noqa: E501
        :rtype: PointAltitude
        """
        return util.deserialize_model(dikt, cls)

    @property
    def shape(self):
        """Gets the shape of this PointAltitude.


        :return: The shape of this PointAltitude.
        :rtype: SupportedGADShapes
        """
        return self._shape

    @shape.setter
    def shape(self, shape):
        """Sets the shape of this PointAltitude.


        :param shape: The shape of this PointAltitude.
        :type shape: SupportedGADShapes
        """
        if shape is None:
            raise ValueError("Invalid value for `shape`, must not be `None`")  # noqa: E501

        self._shape = shape

    @property
    def point(self):
        """Gets the point of this PointAltitude.


        :return: The point of this PointAltitude.
        :rtype: GeographicalCoordinates
        """
        return self._point

    @point.setter
    def point(self, point):
        """Sets the point of this PointAltitude.


        :param point: The point of this PointAltitude.
        :type point: GeographicalCoordinates
        """
        if point is None:
            raise ValueError("Invalid value for `point`, must not be `None`")  # noqa: E501

        self._point = point

    @property
    def altitude(self):
        """Gets the altitude of this PointAltitude.

        Indicates value of altitude.  # noqa: E501

        :return: The altitude of this PointAltitude.
        :rtype: float
        """
        return self._altitude

    @altitude.setter
    def altitude(self, altitude):
        """Sets the altitude of this PointAltitude.

        Indicates value of altitude.  # noqa: E501

        :param altitude: The altitude of this PointAltitude.
        :type altitude: float
        """
        if altitude is None:
            raise ValueError("Invalid value for `altitude`, must not be `None`")  # noqa: E501
        if altitude is not None and altitude > 32767:  # noqa: E501
            raise ValueError("Invalid value for `altitude`, must be a value less than or equal to `32767`")  # noqa: E501
        if altitude is not None and altitude < -32767:  # noqa: E501
            raise ValueError("Invalid value for `altitude`, must be a value greater than or equal to `-32767`")  # noqa: E501

        self._altitude = altitude
