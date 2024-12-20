# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.pdu_session_info import PduSessionInfo
from openapi_server import util

from openapi_server.models.pdu_session_info import PduSessionInfo  # noqa: E501

class PduSessionInformation(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, pdu_sess_id=None, sess_info=None):  # noqa: E501
        """PduSessionInformation - a model defined in OpenAPI

        :param pdu_sess_id: The pdu_sess_id of this PduSessionInformation.  # noqa: E501
        :type pdu_sess_id: int
        :param sess_info: The sess_info of this PduSessionInformation.  # noqa: E501
        :type sess_info: PduSessionInfo
        """
        self.openapi_types = {
            'pdu_sess_id': int,
            'sess_info': PduSessionInfo
        }

        self.attribute_map = {
            'pdu_sess_id': 'pduSessId',
            'sess_info': 'sessInfo'
        }

        self.pdu_sess_id = pdu_sess_id
        self.sess_info = sess_info

    @classmethod
    def from_dict(cls, dikt) -> 'PduSessionInformation':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The PduSessionInformation of this PduSessionInformation.  # noqa: E501
        :rtype: PduSessionInformation
        """
        return util.deserialize_model(dikt, cls)

    @property
    def pdu_sess_id(self):
        """Gets the pdu_sess_id of this PduSessionInformation.

        Unsigned integer identifying a PDU session, within the range 0 to 255, as specified in  clause 11.2.3.1b, bits 1 to 8, of 3GPP TS 24.007. If the PDU Session ID is allocated by the  Core Network for UEs not supporting N1 mode, reserved range 64 to 95 is used. PDU Session ID  within the reserved range is only visible in the Core Network.    # noqa: E501

        :return: The pdu_sess_id of this PduSessionInformation.
        :rtype: int
        """
        return self._pdu_sess_id

    @pdu_sess_id.setter
    def pdu_sess_id(self, pdu_sess_id):
        """Sets the pdu_sess_id of this PduSessionInformation.

        Unsigned integer identifying a PDU session, within the range 0 to 255, as specified in  clause 11.2.3.1b, bits 1 to 8, of 3GPP TS 24.007. If the PDU Session ID is allocated by the  Core Network for UEs not supporting N1 mode, reserved range 64 to 95 is used. PDU Session ID  within the reserved range is only visible in the Core Network.    # noqa: E501

        :param pdu_sess_id: The pdu_sess_id of this PduSessionInformation.
        :type pdu_sess_id: int
        """
        if pdu_sess_id is not None and pdu_sess_id > 255:  # noqa: E501
            raise ValueError("Invalid value for `pdu_sess_id`, must be a value less than or equal to `255`")  # noqa: E501
        if pdu_sess_id is not None and pdu_sess_id < 0:  # noqa: E501
            raise ValueError("Invalid value for `pdu_sess_id`, must be a value greater than or equal to `0`")  # noqa: E501

        self._pdu_sess_id = pdu_sess_id

    @property
    def sess_info(self):
        """Gets the sess_info of this PduSessionInformation.


        :return: The sess_info of this PduSessionInformation.
        :rtype: PduSessionInfo
        """
        return self._sess_info

    @sess_info.setter
    def sess_info(self, sess_info):
        """Sets the sess_info of this PduSessionInformation.


        :param sess_info: The sess_info of this PduSessionInformation.
        :type sess_info: PduSessionInfo
        """

        self._sess_info = sess_info
