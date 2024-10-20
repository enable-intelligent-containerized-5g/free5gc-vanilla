# coding: utf-8

from __future__ import absolute_import
from datetime import date, datetime  # noqa: F401

from typing import List, Dict  # noqa: F401

from openapi_server.models.base_model_ import Model
from openapi_server import util


class MetricsReportingConfiguration(Model):
    """NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).

    Do not edit the class manually.
    """

    def __init__(self, metrics_reporting_configuration_id=None, scheme=None, data_network_name=None, reporting_interval=None, sample_percentage=None, url_filters=None, metrics=None):  # noqa: E501
        """MetricsReportingConfiguration - a model defined in OpenAPI

        :param metrics_reporting_configuration_id: The metrics_reporting_configuration_id of this MetricsReportingConfiguration.  # noqa: E501
        :type metrics_reporting_configuration_id: str
        :param scheme: The scheme of this MetricsReportingConfiguration.  # noqa: E501
        :type scheme: str
        :param data_network_name: The data_network_name of this MetricsReportingConfiguration.  # noqa: E501
        :type data_network_name: str
        :param reporting_interval: The reporting_interval of this MetricsReportingConfiguration.  # noqa: E501
        :type reporting_interval: int
        :param sample_percentage: The sample_percentage of this MetricsReportingConfiguration.  # noqa: E501
        :type sample_percentage: float
        :param url_filters: The url_filters of this MetricsReportingConfiguration.  # noqa: E501
        :type url_filters: List[str]
        :param metrics: The metrics of this MetricsReportingConfiguration.  # noqa: E501
        :type metrics: List[str]
        """
        self.openapi_types = {
            'metrics_reporting_configuration_id': str,
            'scheme': str,
            'data_network_name': str,
            'reporting_interval': int,
            'sample_percentage': float,
            'url_filters': List[str],
            'metrics': List[str]
        }

        self.attribute_map = {
            'metrics_reporting_configuration_id': 'metricsReportingConfigurationId',
            'scheme': 'scheme',
            'data_network_name': 'dataNetworkName',
            'reporting_interval': 'reportingInterval',
            'sample_percentage': 'samplePercentage',
            'url_filters': 'urlFilters',
            'metrics': 'metrics'
        }

        self.metrics_reporting_configuration_id = metrics_reporting_configuration_id
        self.scheme = scheme
        self.data_network_name = data_network_name
        self.reporting_interval = reporting_interval
        self.sample_percentage = sample_percentage
        self.url_filters = url_filters
        self.metrics = metrics

    @classmethod
    def from_dict(cls, dikt) -> 'MetricsReportingConfiguration':
        """Returns the dict as a model

        :param dikt: A dict.
        :type: dict
        :return: The MetricsReportingConfiguration of this MetricsReportingConfiguration.  # noqa: E501
        :rtype: MetricsReportingConfiguration
        """
        return util.deserialize_model(dikt, cls)

    @property
    def metrics_reporting_configuration_id(self):
        """Gets the metrics_reporting_configuration_id of this MetricsReportingConfiguration.

        String chosen by the 5GMS AF to serve as an identifier in a resource URI.  # noqa: E501

        :return: The metrics_reporting_configuration_id of this MetricsReportingConfiguration.
        :rtype: str
        """
        return self._metrics_reporting_configuration_id

    @metrics_reporting_configuration_id.setter
    def metrics_reporting_configuration_id(self, metrics_reporting_configuration_id):
        """Sets the metrics_reporting_configuration_id of this MetricsReportingConfiguration.

        String chosen by the 5GMS AF to serve as an identifier in a resource URI.  # noqa: E501

        :param metrics_reporting_configuration_id: The metrics_reporting_configuration_id of this MetricsReportingConfiguration.
        :type metrics_reporting_configuration_id: str
        """
        if metrics_reporting_configuration_id is None:
            raise ValueError("Invalid value for `metrics_reporting_configuration_id`, must not be `None`")  # noqa: E501

        self._metrics_reporting_configuration_id = metrics_reporting_configuration_id

    @property
    def scheme(self):
        """Gets the scheme of this MetricsReportingConfiguration.

        String providing an URI formatted according to RFC 3986.  # noqa: E501

        :return: The scheme of this MetricsReportingConfiguration.
        :rtype: str
        """
        return self._scheme

    @scheme.setter
    def scheme(self, scheme):
        """Sets the scheme of this MetricsReportingConfiguration.

        String providing an URI formatted according to RFC 3986.  # noqa: E501

        :param scheme: The scheme of this MetricsReportingConfiguration.
        :type scheme: str
        """
        if scheme is None:
            raise ValueError("Invalid value for `scheme`, must not be `None`")  # noqa: E501

        self._scheme = scheme

    @property
    def data_network_name(self):
        """Gets the data_network_name of this MetricsReportingConfiguration.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :return: The data_network_name of this MetricsReportingConfiguration.
        :rtype: str
        """
        return self._data_network_name

    @data_network_name.setter
    def data_network_name(self, data_network_name):
        """Sets the data_network_name of this MetricsReportingConfiguration.

        String representing a Data Network as defined in clause 9A of 3GPP TS 23.003;  it shall contain either a DNN Network Identifier, or a full DNN with both the Network  Identifier and Operator Identifier, as specified in 3GPP TS 23.003 clause 9.1.1 and 9.1.2. It shall be coded as string in which the labels are separated by dots  (e.g. \"Label1.Label2.Label3\").   # noqa: E501

        :param data_network_name: The data_network_name of this MetricsReportingConfiguration.
        :type data_network_name: str
        """

        self._data_network_name = data_network_name

    @property
    def reporting_interval(self):
        """Gets the reporting_interval of this MetricsReportingConfiguration.

        indicating a time in seconds.  # noqa: E501

        :return: The reporting_interval of this MetricsReportingConfiguration.
        :rtype: int
        """
        return self._reporting_interval

    @reporting_interval.setter
    def reporting_interval(self, reporting_interval):
        """Sets the reporting_interval of this MetricsReportingConfiguration.

        indicating a time in seconds.  # noqa: E501

        :param reporting_interval: The reporting_interval of this MetricsReportingConfiguration.
        :type reporting_interval: int
        """

        self._reporting_interval = reporting_interval

    @property
    def sample_percentage(self):
        """Gets the sample_percentage of this MetricsReportingConfiguration.


        :return: The sample_percentage of this MetricsReportingConfiguration.
        :rtype: float
        """
        return self._sample_percentage

    @sample_percentage.setter
    def sample_percentage(self, sample_percentage):
        """Sets the sample_percentage of this MetricsReportingConfiguration.


        :param sample_percentage: The sample_percentage of this MetricsReportingConfiguration.
        :type sample_percentage: float
        """
        if sample_percentage is not None and sample_percentage > 100.0:  # noqa: E501
            raise ValueError("Invalid value for `sample_percentage`, must be a value less than or equal to `100.0`")  # noqa: E501
        if sample_percentage is not None and sample_percentage < 0.0:  # noqa: E501
            raise ValueError("Invalid value for `sample_percentage`, must be a value greater than or equal to `0.0`")  # noqa: E501

        self._sample_percentage = sample_percentage

    @property
    def url_filters(self):
        """Gets the url_filters of this MetricsReportingConfiguration.


        :return: The url_filters of this MetricsReportingConfiguration.
        :rtype: List[str]
        """
        return self._url_filters

    @url_filters.setter
    def url_filters(self, url_filters):
        """Sets the url_filters of this MetricsReportingConfiguration.


        :param url_filters: The url_filters of this MetricsReportingConfiguration.
        :type url_filters: List[str]
        """
        if url_filters is not None and len(url_filters) < 1:
            raise ValueError("Invalid value for `url_filters`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._url_filters = url_filters

    @property
    def metrics(self):
        """Gets the metrics of this MetricsReportingConfiguration.


        :return: The metrics of this MetricsReportingConfiguration.
        :rtype: List[str]
        """
        return self._metrics

    @metrics.setter
    def metrics(self, metrics):
        """Sets the metrics of this MetricsReportingConfiguration.


        :param metrics: The metrics of this MetricsReportingConfiguration.
        :type metrics: List[str]
        """
        if metrics is not None and len(metrics) < 1:
            raise ValueError("Invalid value for `metrics`, number of items must be greater than or equal to `1`")  # noqa: E501

        self._metrics = metrics
