# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server.models.ml_model_addr import MLModelAddr
from openapi_server import util

from openapi_server.models.ml_model_addr import MLModelAddr  # noqa: E501

class MLModelInfo(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, ml_file_addrs=None, model_prov_id=None, model_prov_set_id=None):  # noqa: E501
        """MLModelInfo - a model defined in OpenAPI

        :param ml_file_addrs: The ml_file_addrs of this MLModelInfo.  # noqa: E501
        :type ml_file_addrs: List[MLModelAddr]
        :param model_prov_id: The model_prov_id of this MLModelInfo.  # noqa: E501
        :type model_prov_id: str
        :param model_prov_set_id: The model_prov_set_id of this MLModelInfo.  # noqa: E501
        :type model_prov_set_id: str
        """
        self.openapi_types = {
            'ml_file_addrs': List[MLModelAddr],
            'model_prov_id': str,
            'model_prov_set_id': str
        }

        self.attribute_map = {
            'ml_file_addrs': 'mlFileAddrs',
            'model_prov_id': 'modelProvId',
            'model_prov_set_id': 'modelProvSetId'
        }

        self.ml_file_addrs = ml_file_addrs
        self.model_prov_id = model_prov_id
        self.model_prov_set_id = model_prov_set_id

    @classmethod
    def from_dict(cls, dikt) -> 'MLModelInfo':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The MLModelInfo of this MLModelInfo.  # noqa: E501
        :rtype: MLModelInfo
        """
        return util.deserialize_model(dikt, cls)

    @property
    def ml_file_addrs(self):
        """Gets the ml_file_addrs of this MLModelInfo.


        :return: The ml_file_addrs of this MLModelInfo.
        :rtype: List[MLModelAddr]
        """
        return self._ml_file_addrs

    @ml_file_addrs.setter
    def ml_file_addrs(self, ml_file_addrs):
        """Sets the ml_file_addrs of this MLModelInfo.


        :param ml_file_addrs: The ml_file_addrs of this MLModelInfo.
        :type ml_file_addrs: List[MLModelAddr]
        """
        if ml_file_addrs is not None and len(ml_file_addrs) < 1:
            raise ValueError("Invalid value for `ml_file_addrs`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._ml_file_addrs = ml_file_addrs

    @property
    def model_prov_id(self):
        """Gets the model_prov_id of this MLModelInfo.

        String uniquely identifying a NF instance. The format of the NF Instance ID shall be a  Universally Unique Identifier (UUID) version 4, as described in IETF RFC 4122.    # noqa: E501

        :return: The model_prov_id of this MLModelInfo.
        :rtype: str
        """
        return self._model_prov_id

    @model_prov_id.setter
    def model_prov_id(self, model_prov_id):
        """Sets the model_prov_id of this MLModelInfo.

        String uniquely identifying a NF instance. The format of the NF Instance ID shall be a  Universally Unique Identifier (UUID) version 4, as described in IETF RFC 4122.    # noqa: E501

        :param model_prov_id: The model_prov_id of this MLModelInfo.
        :type model_prov_id: str
        """

        self._model_prov_id = model_prov_id

    @property
    def model_prov_set_id(self):
        """Gets the model_prov_set_id of this MLModelInfo.

        NF Set Identifier (see clause 28.12 of 3GPP TS 23.003), formatted as the following string \"set<Set ID>.<nftype>set.5gc.mnc<MNC>.mcc<MCC>\", or  \"set<SetID>.<NFType>set.5gc.nid<NID>.mnc<MNC>.mcc<MCC>\" with  <MCC> encoded as defined in clause 5.4.2 (\"Mcc\" data type definition)  <MNC> encoding the Mobile Network Code part of the PLMN, comprising 3 digits.    If there are only 2 significant digits in the MNC, one \"0\" digit shall be inserted    at the left side to fill the 3 digits coding of MNC.  Pattern: '^[0-9]{3}$' <NFType> encoded as a value defined in Table 6.1.6.3.3-1 of 3GPP TS 29.510 but    with lower case characters <Set ID> encoded as a string of characters consisting of    alphabetic characters (A-Z and a-z), digits (0-9) and/or the hyphen (-) and that    shall end with either an alphabetic character or a digit.    # noqa: E501

        :return: The model_prov_set_id of this MLModelInfo.
        :rtype: str
        """
        return self._model_prov_set_id

    @model_prov_set_id.setter
    def model_prov_set_id(self, model_prov_set_id):
        """Sets the model_prov_set_id of this MLModelInfo.

        NF Set Identifier (see clause 28.12 of 3GPP TS 23.003), formatted as the following string \"set<Set ID>.<nftype>set.5gc.mnc<MNC>.mcc<MCC>\", or  \"set<SetID>.<NFType>set.5gc.nid<NID>.mnc<MNC>.mcc<MCC>\" with  <MCC> encoded as defined in clause 5.4.2 (\"Mcc\" data type definition)  <MNC> encoding the Mobile Network Code part of the PLMN, comprising 3 digits.    If there are only 2 significant digits in the MNC, one \"0\" digit shall be inserted    at the left side to fill the 3 digits coding of MNC.  Pattern: '^[0-9]{3}$' <NFType> encoded as a value defined in Table 6.1.6.3.3-1 of 3GPP TS 29.510 but    with lower case characters <Set ID> encoded as a string of characters consisting of    alphabetic characters (A-Z and a-z), digits (0-9) and/or the hyphen (-) and that    shall end with either an alphabetic character or a digit.    # noqa: E501

        :param model_prov_set_id: The model_prov_set_id of this MLModelInfo.
        :type model_prov_set_id: str
        """

        self._model_prov_set_id = model_prov_set_id
