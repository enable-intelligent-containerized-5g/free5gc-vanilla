# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class DispersionOrderingCriterionAnyOf(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    """
    allowed enum values
    """
    TIME_SLOT_START = "TIME_SLOT_START"
    DISPERSION = "DISPERSION"
    CLASSIFICATION = "CLASSIFICATION"
    RANKING = "RANKING"
    PERCENTILE_RANKING = "PERCENTILE_RANKING"
    def __init__(self):  # noqa: E501
        """DispersionOrderingCriterionAnyOf - a model defined in OpenAPI

        """
        self.openapi_types = {
        }

        self.attribute_map = {
        }

    @classmethod
    def from_dict(cls, dikt) -> 'DispersionOrderingCriterionAnyOf':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The DispersionOrderingCriterion_anyOf of this DispersionOrderingCriterionAnyOf.  # noqa: E501
        :rtype: DispersionOrderingCriterionAnyOf
        """
        return util.deserialize_model(dikt, cls)