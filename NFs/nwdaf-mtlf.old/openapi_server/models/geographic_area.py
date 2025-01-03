# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ellipsoid_arc import EllipsoidArc
from openapi_server.models.geographical_coordinates import GeographicalCoordinates
from openapi_server.models.point import Point
from openapi_server.models.point_altitude import PointAltitude
from openapi_server.models.point_altitude_uncertainty import PointAltitudeUncertainty
from openapi_server.models.point_uncertainty_circle import PointUncertaintyCircle
from openapi_server.models.point_uncertainty_ellipse import PointUncertaintyEllipse
from openapi_server.models.polygon import Polygon
from openapi_server.models.supported_gad_shapes import SupportedGADShapes
from openapi_server.models.uncertainty_ellipse import UncertaintyEllipse
from openapi_server import util

from openapi_server.models.ellipsoid_arc import EllipsoidArc  # noqa: E501
from openapi_server.models.geographical_coordinates import GeographicalCoordinates  # noqa: E501
from openapi_server.models.point import Point  # noqa: E501
from openapi_server.models.point_altitude import PointAltitude  # noqa: E501
from openapi_server.models.point_altitude_uncertainty import PointAltitudeUncertainty  # noqa: E501
from openapi_server.models.point_uncertainty_circle import PointUncertaintyCircle  # noqa: E501
from openapi_server.models.point_uncertainty_ellipse import PointUncertaintyEllipse  # noqa: E501
from openapi_server.models.polygon import Polygon  # noqa: E501
from openapi_server.models.supported_gad_shapes import SupportedGADShapes  # noqa: E501
from openapi_server.models.uncertainty_ellipse import UncertaintyEllipse  # noqa: E501

class GeographicArea(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, shape=None, point=None, uncertainty=None, uncertainty_ellipse=None, confidence=None, point_list=None, altitude=None, uncertainty_altitude=None, inner_radius=None, uncertainty_radius=None, offset_angle=None, included_angle=None):  # noqa: E501
        """GeographicArea - a model defined in OpenAPI

        :param shape: The shape of this GeographicArea.  # noqa: E501
        :type shape: SupportedGADShapes
        :param point: The point of this GeographicArea.  # noqa: E501
        :type point: GeographicalCoordinates
        :param uncertainty: The uncertainty of this GeographicArea.  # noqa: E501
        :type uncertainty: float
        :param uncertainty_ellipse: The uncertainty_ellipse of this GeographicArea.  # noqa: E501
        :type uncertainty_ellipse: UncertaintyEllipse
        :param confidence: The confidence of this GeographicArea.  # noqa: E501
        :type confidence: int
        :param point_list: The point_list of this GeographicArea.  # noqa: E501
        :type point_list: List[GeographicalCoordinates]
        :param altitude: The altitude of this GeographicArea.  # noqa: E501
        :type altitude: float
        :param uncertainty_altitude: The uncertainty_altitude of this GeographicArea.  # noqa: E501
        :type uncertainty_altitude: float
        :param inner_radius: The inner_radius of this GeographicArea.  # noqa: E501
        :type inner_radius: int
        :param uncertainty_radius: The uncertainty_radius of this GeographicArea.  # noqa: E501
        :type uncertainty_radius: float
        :param offset_angle: The offset_angle of this GeographicArea.  # noqa: E501
        :type offset_angle: int
        :param included_angle: The included_angle of this GeographicArea.  # noqa: E501
        :type included_angle: int
        """
        self.openapi_types = {
            'shape': SupportedGADShapes,
            'point': GeographicalCoordinates,
            'uncertainty': float,
            'uncertainty_ellipse': UncertaintyEllipse,
            'confidence': int,
            'point_list': List[GeographicalCoordinates],
            'altitude': float,
            'uncertainty_altitude': float,
            'inner_radius': int,
            'uncertainty_radius': float,
            'offset_angle': int,
            'included_angle': int
        }

        self.attribute_map = {
            'shape': 'shape',
            'point': 'point',
            'uncertainty': 'uncertainty',
            'uncertainty_ellipse': 'uncertaintyEllipse',
            'confidence': 'confidence',
            'point_list': 'pointList',
            'altitude': 'altitude',
            'uncertainty_altitude': 'uncertaintyAltitude',
            'inner_radius': 'innerRadius',
            'uncertainty_radius': 'uncertaintyRadius',
            'offset_angle': 'offsetAngle',
            'included_angle': 'includedAngle'
        }

        self.shape = shape
        self.point = point
        self.uncertainty = uncertainty
        self.uncertainty_ellipse = uncertainty_ellipse
        self.confidence = confidence
        self.point_list = point_list
        self.altitude = altitude
        self.uncertainty_altitude = uncertainty_altitude
        self.inner_radius = inner_radius
        self.uncertainty_radius = uncertainty_radius
        self.offset_angle = offset_angle
        self.included_angle = included_angle

    @classmethod
    def from_dict(cls, dikt) -> 'GeographicArea':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The GeographicArea of this GeographicArea.  # noqa: E501
        :rtype: GeographicArea
        """
        return util.deserialize_model(dikt, cls)

    @property
    def shape(self):
        """Gets the shape of this GeographicArea.


        :return: The shape of this GeographicArea.
        :rtype: SupportedGADShapes
        """
        return self._shape

    @shape.setter
    def shape(self, shape):
        """Sets the shape of this GeographicArea.


        :param shape: The shape of this GeographicArea.
        :type shape: SupportedGADShapes
        """
        if shape is None:
            raise ValueError("Invalid value for `shape`, must not be `None`")  # noqa: E501

        self._shape = shape

    @property
    def point(self):
        """Gets the point of this GeographicArea.


        :return: The point of this GeographicArea.
        :rtype: GeographicalCoordinates
        """
        return self._point

    @point.setter
    def point(self, point):
        """Sets the point of this GeographicArea.


        :param point: The point of this GeographicArea.
        :type point: GeographicalCoordinates
        """
        if point is None:
            raise ValueError("Invalid value for `point`, must not be `None`")  # noqa: E501

        self._point = point

    @property
    def uncertainty(self):
        """Gets the uncertainty of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :return: The uncertainty of this GeographicArea.
        :rtype: float
        """
        return self._uncertainty

    @uncertainty.setter
    def uncertainty(self, uncertainty):
        """Sets the uncertainty of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :param uncertainty: The uncertainty of this GeographicArea.
        :type uncertainty: float
        """
        if uncertainty is None:
            raise ValueError("Invalid value for `uncertainty`, must not be `None`")  # noqa: E501
        if uncertainty is not None and uncertainty < 0:  # noqa: E501
            raise ValueError("Invalid value for `uncertainty`, must be a value greater than or equal to `0`")  # noqa: E501

        self._uncertainty = uncertainty

    @property
    def uncertainty_ellipse(self):
        """Gets the uncertainty_ellipse of this GeographicArea.


        :return: The uncertainty_ellipse of this GeographicArea.
        :rtype: UncertaintyEllipse
        """
        return self._uncertainty_ellipse

    @uncertainty_ellipse.setter
    def uncertainty_ellipse(self, uncertainty_ellipse):
        """Sets the uncertainty_ellipse of this GeographicArea.


        :param uncertainty_ellipse: The uncertainty_ellipse of this GeographicArea.
        :type uncertainty_ellipse: UncertaintyEllipse
        """
        if uncertainty_ellipse is None:
            raise ValueError("Invalid value for `uncertainty_ellipse`, must not be `None`")  # noqa: E501

        self._uncertainty_ellipse = uncertainty_ellipse

    @property
    def confidence(self):
        """Gets the confidence of this GeographicArea.

        Indicates value of confidence.  # noqa: E501

        :return: The confidence of this GeographicArea.
        :rtype: int
        """
        return self._confidence

    @confidence.setter
    def confidence(self, confidence):
        """Sets the confidence of this GeographicArea.

        Indicates value of confidence.  # noqa: E501

        :param confidence: The confidence of this GeographicArea.
        :type confidence: int
        """
        if confidence is None:
            raise ValueError("Invalid value for `confidence`, must not be `None`")  # noqa: E501
        if confidence is not None and confidence > 100:  # noqa: E501
            raise ValueError("Invalid value for `confidence`, must be a value less than or equal to `100`")  # noqa: E501
        if confidence is not None and confidence < 0:  # noqa: E501
            raise ValueError("Invalid value for `confidence`, must be a value greater than or equal to `0`")  # noqa: E501

        self._confidence = confidence

    @property
    def point_list(self):
        """Gets the point_list of this GeographicArea.

        List of points.  # noqa: E501

        :return: The point_list of this GeographicArea.
        :rtype: List[GeographicalCoordinates]
        """
        return self._point_list

    @point_list.setter
    def point_list(self, point_list):
        """Sets the point_list of this GeographicArea.

        List of points.  # noqa: E501

        :param point_list: The point_list of this GeographicArea.
        :type point_list: List[GeographicalCoordinates]
        """
        if point_list is None:
            raise ValueError("Invalid value for `point_list`, must not be `None`")  # noqa: E501
        if point_list is not None and len(point_list) > 15:
            raise ValueError("Invalid value for `point_list`, number of items must be less than or equal to `15`")  # noqa: E501
        if point_list is not None and len(point_list) < 3:
            raise ValueError("Invalid value for `point_list`, number of items must be greater than or equal to `3`")  # noqa: E501

        self._point_list = point_list

    @property
    def altitude(self):
        """Gets the altitude of this GeographicArea.

        Indicates value of altitude.  # noqa: E501

        :return: The altitude of this GeographicArea.
        :rtype: float
        """
        return self._altitude

    @altitude.setter
    def altitude(self, altitude):
        """Sets the altitude of this GeographicArea.

        Indicates value of altitude.  # noqa: E501

        :param altitude: The altitude of this GeographicArea.
        :type altitude: float
        """
        if altitude is None:
            raise ValueError("Invalid value for `altitude`, must not be `None`")  # noqa: E501
        if altitude is not None and altitude > 32767:  # noqa: E501
            raise ValueError("Invalid value for `altitude`, must be a value less than or equal to `32767`")  # noqa: E501
        if altitude is not None and altitude < -32767:  # noqa: E501
            raise ValueError("Invalid value for `altitude`, must be a value greater than or equal to `-32767`")  # noqa: E501

        self._altitude = altitude

    @property
    def uncertainty_altitude(self):
        """Gets the uncertainty_altitude of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :return: The uncertainty_altitude of this GeographicArea.
        :rtype: float
        """
        return self._uncertainty_altitude

    @uncertainty_altitude.setter
    def uncertainty_altitude(self, uncertainty_altitude):
        """Sets the uncertainty_altitude of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :param uncertainty_altitude: The uncertainty_altitude of this GeographicArea.
        :type uncertainty_altitude: float
        """
        if uncertainty_altitude is None:
            raise ValueError("Invalid value for `uncertainty_altitude`, must not be `None`")  # noqa: E501
        if uncertainty_altitude is not None and uncertainty_altitude < 0:  # noqa: E501
            raise ValueError("Invalid value for `uncertainty_altitude`, must be a value greater than or equal to `0`")  # noqa: E501

        self._uncertainty_altitude = uncertainty_altitude

    @property
    def inner_radius(self):
        """Gets the inner_radius of this GeographicArea.

        Indicates value of the inner radius.  # noqa: E501

        :return: The inner_radius of this GeographicArea.
        :rtype: int
        """
        return self._inner_radius

    @inner_radius.setter
    def inner_radius(self, inner_radius):
        """Sets the inner_radius of this GeographicArea.

        Indicates value of the inner radius.  # noqa: E501

        :param inner_radius: The inner_radius of this GeographicArea.
        :type inner_radius: int
        """
        if inner_radius is None:
            raise ValueError("Invalid value for `inner_radius`, must not be `None`")  # noqa: E501
        if inner_radius is not None and inner_radius > 327675:  # noqa: E501
            raise ValueError("Invalid value for `inner_radius`, must be a value less than or equal to `327675`")  # noqa: E501
        if inner_radius is not None and inner_radius < 0:  # noqa: E501
            raise ValueError("Invalid value for `inner_radius`, must be a value greater than or equal to `0`")  # noqa: E501

        self._inner_radius = inner_radius

    @property
    def uncertainty_radius(self):
        """Gets the uncertainty_radius of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :return: The uncertainty_radius of this GeographicArea.
        :rtype: float
        """
        return self._uncertainty_radius

    @uncertainty_radius.setter
    def uncertainty_radius(self, uncertainty_radius):
        """Sets the uncertainty_radius of this GeographicArea.

        Indicates value of uncertainty.  # noqa: E501

        :param uncertainty_radius: The uncertainty_radius of this GeographicArea.
        :type uncertainty_radius: float
        """
        if uncertainty_radius is None:
            raise ValueError("Invalid value for `uncertainty_radius`, must not be `None`")  # noqa: E501
        if uncertainty_radius is not None and uncertainty_radius < 0:  # noqa: E501
            raise ValueError("Invalid value for `uncertainty_radius`, must be a value greater than or equal to `0`")  # noqa: E501

        self._uncertainty_radius = uncertainty_radius

    @property
    def offset_angle(self):
        """Gets the offset_angle of this GeographicArea.

        Indicates value of angle.  # noqa: E501

        :return: The offset_angle of this GeographicArea.
        :rtype: int
        """
        return self._offset_angle

    @offset_angle.setter
    def offset_angle(self, offset_angle):
        """Sets the offset_angle of this GeographicArea.

        Indicates value of angle.  # noqa: E501

        :param offset_angle: The offset_angle of this GeographicArea.
        :type offset_angle: int
        """
        if offset_angle is None:
            raise ValueError("Invalid value for `offset_angle`, must not be `None`")  # noqa: E501
        if offset_angle is not None and offset_angle > 360:  # noqa: E501
            raise ValueError("Invalid value for `offset_angle`, must be a value less than or equal to `360`")  # noqa: E501
        if offset_angle is not None and offset_angle < 0:  # noqa: E501
            raise ValueError("Invalid value for `offset_angle`, must be a value greater than or equal to `0`")  # noqa: E501

        self._offset_angle = offset_angle

    @property
    def included_angle(self):
        """Gets the included_angle of this GeographicArea.

        Indicates value of angle.  # noqa: E501

        :return: The included_angle of this GeographicArea.
        :rtype: int
        """
        return self._included_angle

    @included_angle.setter
    def included_angle(self, included_angle):
        """Sets the included_angle of this GeographicArea.

        Indicates value of angle.  # noqa: E501

        :param included_angle: The included_angle of this GeographicArea.
        :type included_angle: int
        """
        if included_angle is None:
            raise ValueError("Invalid value for `included_angle`, must not be `None`")  # noqa: E501
        if included_angle is not None and included_angle > 360:  # noqa: E501
            raise ValueError("Invalid value for `included_angle`, must be a value less than or equal to `360`")  # noqa: E501
        if included_angle is not None and included_angle < 0:  # noqa: E501
            raise ValueError("Invalid value for `included_angle`, must be a value greater than or equal to `0`")  # noqa: E501

        self._included_angle = included_angle
