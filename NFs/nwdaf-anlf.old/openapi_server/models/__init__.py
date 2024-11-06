# coding: utf-8

# flake8: noqa
from __future__ import absolute_import
# import models into model package
from openapi_server.models.abnormal_behaviour import AbnormalBehaviour
from openapi_server.models.access_state_transition_type import AccessStateTransitionType
from openapi_server.models.access_state_transition_type_any_of import AccessStateTransitionTypeAnyOf
from openapi_server.models.access_token_err import AccessTokenErr
from openapi_server.models.access_token_req import AccessTokenReq
from openapi_server.models.access_type import AccessType
from openapi_server.models.accuracy import Accuracy
from openapi_server.models.accuracy_any_of import AccuracyAnyOf
from openapi_server.models.addition_info_analytics_info_request import AdditionInfoAnalyticsInfoRequest
from openapi_server.models.additional_measurement import AdditionalMeasurement
from openapi_server.models.addr_fqdn import AddrFqdn
from openapi_server.models.address_list import AddressList
from openapi_server.models.adrf_data_type import AdrfDataType
from openapi_server.models.adrf_data_type_any_of import AdrfDataTypeAnyOf
from openapi_server.models.af_event import AfEvent
from openapi_server.models.af_event_exposure_notif import AfEventExposureNotif
from openapi_server.models.af_event_exposure_subsc import AfEventExposureSubsc
from openapi_server.models.af_event_notification import AfEventNotification
from openapi_server.models.amf_event import AmfEvent
from openapi_server.models.amf_event_area import AmfEventArea
from openapi_server.models.amf_event_mode import AmfEventMode
from openapi_server.models.amf_event_notification import AmfEventNotification
from openapi_server.models.amf_event_report import AmfEventReport
from openapi_server.models.amf_event_state import AmfEventState
from openapi_server.models.amf_event_subs_sync_info import AmfEventSubsSyncInfo
from openapi_server.models.amf_event_subscription import AmfEventSubscription
from openapi_server.models.amf_event_subscription_info import AmfEventSubscriptionInfo
from openapi_server.models.amf_event_trigger import AmfEventTrigger
from openapi_server.models.amf_event_trigger_any_of import AmfEventTriggerAnyOf
from openapi_server.models.amf_event_type import AmfEventType
from openapi_server.models.amf_event_type_any_of import AmfEventTypeAnyOf
from openapi_server.models.analytics_context_identifier import AnalyticsContextIdentifier
from openapi_server.models.analytics_data import AnalyticsData
from openapi_server.models.analytics_metadata import AnalyticsMetadata
from openapi_server.models.analytics_metadata_any_of import AnalyticsMetadataAnyOf
from openapi_server.models.analytics_metadata_indication import AnalyticsMetadataIndication
from openapi_server.models.analytics_metadata_info import AnalyticsMetadataInfo
from openapi_server.models.analytics_subset import AnalyticsSubset
from openapi_server.models.analytics_subset_any_of import AnalyticsSubsetAnyOf
from openapi_server.models.app_list_for_ue_comm import AppListForUeComm
from openapi_server.models.application_volume import ApplicationVolume
from openapi_server.models.applied_smcc_type import AppliedSmccType
from openapi_server.models.applied_smcc_type_any_of import AppliedSmccTypeAnyOf
from openapi_server.models.association_type import AssociationType
from openapi_server.models.association_type_any_of import AssociationTypeAnyOf
from openapi_server.models.battery_indication import BatteryIndication
from openapi_server.models.bw_requirement import BwRequirement
from openapi_server.models.cell_global_id import CellGlobalId
from openapi_server.models.change_of_supi_pei_association_report import ChangeOfSupiPeiAssociationReport
from openapi_server.models.charg_policy_invocation_collection import ChargPolicyInvocationCollection
from openapi_server.models.circumstance_description import CircumstanceDescription
from openapi_server.models.civic_address import CivicAddress
from openapi_server.models.class_criterion import ClassCriterion
from openapi_server.models.cm_info import CmInfo
from openapi_server.models.cm_info_report import CmInfoReport
from openapi_server.models.cm_state import CmState
from openapi_server.models.cm_state_any_of import CmStateAnyOf
from openapi_server.models.cn_type import CnType
from openapi_server.models.cn_type_any_of import CnTypeAnyOf
from openapi_server.models.cn_type_change_report import CnTypeChangeReport
from openapi_server.models.collective_behaviour_filter import CollectiveBehaviourFilter
from openapi_server.models.collective_behaviour_filter_type import CollectiveBehaviourFilterType
from openapi_server.models.collective_behaviour_filter_type_any_of import CollectiveBehaviourFilterTypeAnyOf
from openapi_server.models.collective_behaviour_info import CollectiveBehaviourInfo
from openapi_server.models.communication_collection import CommunicationCollection
from openapi_server.models.communication_failure import CommunicationFailure
from openapi_server.models.congestion_info import CongestionInfo
from openapi_server.models.congestion_type import CongestionType
from openapi_server.models.congestion_type_any_of import CongestionTypeAnyOf
from openapi_server.models.consumer_nf_information import ConsumerNfInformation
from openapi_server.models.consumption_collection import ConsumptionCollection
from openapi_server.models.context_data import ContextData
from openapi_server.models.context_element import ContextElement
from openapi_server.models.context_id_list import ContextIdList
from openapi_server.models.context_info import ContextInfo
from openapi_server.models.context_type import ContextType
from openapi_server.models.context_type_any_of import ContextTypeAnyOf
from openapi_server.models.data_notification import DataNotification
from openapi_server.models.data_subscription import DataSubscription
from openapi_server.models.datalink_reporting_configuration import DatalinkReportingConfiguration
from openapi_server.models.dataset_statistical_property import DatasetStatisticalProperty
from openapi_server.models.dataset_statistical_property_any_of import DatasetStatisticalPropertyAnyOf
from openapi_server.models.ddd_traffic_descriptor import DddTrafficDescriptor
from openapi_server.models.dispersion_area import DispersionArea
from openapi_server.models.dispersion_class import DispersionClass
from openapi_server.models.dispersion_class_one_of import DispersionClassOneOf
from openapi_server.models.dispersion_collection import DispersionCollection
from openapi_server.models.dispersion_collection1 import DispersionCollection1
from openapi_server.models.dispersion_info import DispersionInfo
from openapi_server.models.dispersion_ordering_criterion import DispersionOrderingCriterion
from openapi_server.models.dispersion_ordering_criterion_any_of import DispersionOrderingCriterionAnyOf
from openapi_server.models.dispersion_requirement import DispersionRequirement
from openapi_server.models.dispersion_type import DispersionType
from openapi_server.models.dispersion_type_one_of import DispersionTypeOneOf
from openapi_server.models.dl_data_delivery_status import DlDataDeliveryStatus
from openapi_server.models.dl_data_delivery_status_any_of import DlDataDeliveryStatusAnyOf
from openapi_server.models.dn_perf import DnPerf
from openapi_server.models.dn_perf_info import DnPerfInfo
from openapi_server.models.dn_perf_ordering_criterion import DnPerfOrderingCriterion
from openapi_server.models.dn_perf_ordering_criterion_any_of import DnPerfOrderingCriterionAnyOf
from openapi_server.models.dn_performance_req import DnPerformanceReq
from openapi_server.models.dnai_change_type import DnaiChangeType
from openapi_server.models.dnai_change_type_any_of import DnaiChangeTypeAnyOf
from openapi_server.models.ecgi import Ecgi
from openapi_server.models.ee_subscription import EeSubscription
from openapi_server.models.ellipsoid_arc import EllipsoidArc
from openapi_server.models.ellipsoid_arc_all_of import EllipsoidArcAllOf
from openapi_server.models.eth_flow_description import EthFlowDescription
from openapi_server.models.eutra_location import EutraLocation
from openapi_server.models.event_filter import EventFilter
from openapi_server.models.event_filter1 import EventFilter1
from openapi_server.models.event_id import EventId
from openapi_server.models.event_id_any_of import EventIdAnyOf
from openapi_server.models.event_notification import EventNotification
from openapi_server.models.event_notification1 import EventNotification1
from openapi_server.models.event_report_mode import EventReportMode
from openapi_server.models.event_report_mode_any_of import EventReportModeAnyOf
from openapi_server.models.event_reporting_requirement import EventReportingRequirement
from openapi_server.models.event_subscription import EventSubscription
from openapi_server.models.event_subscription1 import EventSubscription1
from openapi_server.models.event_type import EventType
from openapi_server.models.event_type_any_of import EventTypeAnyOf
from openapi_server.models.events_subs import EventsSubs
from openapi_server.models.exception import Exception
from openapi_server.models.exception_id import ExceptionId
from openapi_server.models.exception_id_any_of import ExceptionIdAnyOf
from openapi_server.models.exception_info import ExceptionInfo
from openapi_server.models.exception_trend import ExceptionTrend
from openapi_server.models.exception_trend_any_of import ExceptionTrendAnyOf
from openapi_server.models.expected_analytics_type import ExpectedAnalyticsType
from openapi_server.models.expected_analytics_type_any_of import ExpectedAnalyticsTypeAnyOf
from openapi_server.models.expected_ue_behaviour_data import ExpectedUeBehaviourData
from openapi_server.models.ext_snssai import ExtSnssai
from openapi_server.models.failure_event_info import FailureEventInfo
from openapi_server.models.flow_direction import FlowDirection
from openapi_server.models.flow_direction_any_of import FlowDirectionAnyOf
from openapi_server.models.flow_info import FlowInfo
from openapi_server.models.gad_shape import GADShape
from openapi_server.models.gnb_id import GNbId
from openapi_server.models.geographic_area import GeographicArea
from openapi_server.models.geographical_coordinates import GeographicalCoordinates
from openapi_server.models.gera_location import GeraLocation
from openapi_server.models.global_ran_node_id import GlobalRanNodeId
from openapi_server.models.guami import Guami
from openapi_server.models.hfc_node_id import HfcNodeId
from openapi_server.models.historical_data import HistoricalData
from openapi_server.models.idle_status_indication import IdleStatusIndication
from openapi_server.models.invalid_param import InvalidParam
from openapi_server.models.ip_addr import IpAddr
from openapi_server.models.ip_eth_flow_description import IpEthFlowDescription
from openapi_server.models.ipv6_addr import Ipv6Addr
from openapi_server.models.ipv6_prefix import Ipv6Prefix
from openapi_server.models.ladn_info import LadnInfo
from openapi_server.models.line_type import LineType
from openapi_server.models.line_type_any_of import LineTypeAnyOf
from openapi_server.models.local2d_point_uncertainty_ellipse import Local2dPointUncertaintyEllipse
from openapi_server.models.local2d_point_uncertainty_ellipse_all_of import Local2dPointUncertaintyEllipseAllOf
from openapi_server.models.local3d_point_uncertainty_ellipsoid import Local3dPointUncertaintyEllipsoid
from openapi_server.models.local3d_point_uncertainty_ellipsoid_all_of import Local3dPointUncertaintyEllipsoidAllOf
from openapi_server.models.local_origin import LocalOrigin
from openapi_server.models.location_accuracy import LocationAccuracy
from openapi_server.models.location_accuracy_any_of import LocationAccuracyAnyOf
from openapi_server.models.location_area import LocationArea
from openapi_server.models.location_area5_g import LocationArea5G
from openapi_server.models.location_area_id import LocationAreaId
from openapi_server.models.location_filter import LocationFilter
from openapi_server.models.location_filter_any_of import LocationFilterAnyOf
from openapi_server.models.location_info import LocationInfo
from openapi_server.models.location_report import LocationReport
from openapi_server.models.location_reporting_configuration import LocationReportingConfiguration
from openapi_server.models.loss_connectivity_cfg import LossConnectivityCfg
from openapi_server.models.loss_connectivity_report import LossConnectivityReport
from openapi_server.models.loss_of_connectivity_reason import LossOfConnectivityReason
from openapi_server.models.loss_of_connectivity_reason_any_of import LossOfConnectivityReasonAnyOf
from openapi_server.models.ml_model_addr import MLModelAddr
from openapi_server.models.ml_model_info import MLModelInfo
from openapi_server.models.ms_access_activity_collection import MSAccessActivityCollection
from openapi_server.models.matching_direction import MatchingDirection
from openapi_server.models.matching_direction_any_of import MatchingDirectionAnyOf
from openapi_server.models.metrics_reporting_configuration import MetricsReportingConfiguration
from openapi_server.models.mm_transaction_location_report_item import MmTransactionLocationReportItem
from openapi_server.models.mm_transaction_slice_report_item import MmTransactionSliceReportItem
from openapi_server.models.model5_gs_user_state import Model5GsUserState
from openapi_server.models.model5_gs_user_state_any_of import Model5GsUserStateAnyOf
from openapi_server.models.model5_gs_user_state_info import Model5GsUserStateInfo
from openapi_server.models.model_info import ModelInfo
from openapi_server.models.monitoring_configuration import MonitoringConfiguration
from openapi_server.models.monitoring_report import MonitoringReport
from openapi_server.models.n3ga_location import N3gaLocation
from openapi_server.models.nf_type import NFType
from openapi_server.models.nf_type_any_of import NFTypeAnyOf
from openapi_server.models.ncgi import Ncgi
from openapi_server.models.nef_event import NefEvent
from openapi_server.models.nef_event_any_of import NefEventAnyOf
from openapi_server.models.nef_event_exposure_notif import NefEventExposureNotif
from openapi_server.models.nef_event_exposure_subsc import NefEventExposureSubsc
from openapi_server.models.nef_event_filter import NefEventFilter
from openapi_server.models.nef_event_notification import NefEventNotification
from openapi_server.models.nef_event_subs import NefEventSubs
from openapi_server.models.net_ass_invocation_collection import NetAssInvocationCollection
from openapi_server.models.network_area_info import NetworkAreaInfo
from openapi_server.models.network_area_info1 import NetworkAreaInfo1
from openapi_server.models.network_perf_info import NetworkPerfInfo
from openapi_server.models.network_perf_requirement import NetworkPerfRequirement
from openapi_server.models.network_perf_type import NetworkPerfType
from openapi_server.models.network_perf_type_any_of import NetworkPerfTypeAnyOf
from openapi_server.models.nf_load_level_information import NfLoadLevelInformation
from openapi_server.models.nf_status import NfStatus
from openapi_server.models.ng_ap_cause import NgApCause
from openapi_server.models.nnwdaf_events_subscription import NnwdafEventsSubscription
from openapi_server.models.notification_flag import NotificationFlag
from openapi_server.models.notification_flag_any_of import NotificationFlagAnyOf
from openapi_server.models.notification_method import NotificationMethod
from openapi_server.models.notification_method1 import NotificationMethod1
from openapi_server.models.notification_method1_any_of import NotificationMethod1AnyOf
from openapi_server.models.notification_method_any_of import NotificationMethodAnyOf
from openapi_server.models.nr_location import NrLocation
from openapi_server.models.nsi_id_info import NsiIdInfo
from openapi_server.models.nsi_load_level_info import NsiLoadLevelInfo
from openapi_server.models.nsmf_event_exposure import NsmfEventExposure
from openapi_server.models.nsmf_event_exposure_notification import NsmfEventExposureNotification
from openapi_server.models.number_average import NumberAverage
from openapi_server.models.nwdaf_event import NwdafEvent
from openapi_server.models.nwdaf_event_any_of import NwdafEventAnyOf
from openapi_server.models.nwdaf_failure_code import NwdafFailureCode
from openapi_server.models.nwdaf_failure_code_any_of import NwdafFailureCodeAnyOf
from openapi_server.models.output_strategy import OutputStrategy
from openapi_server.models.output_strategy_any_of import OutputStrategyAnyOf
from openapi_server.models.partitioning_criteria import PartitioningCriteria
from openapi_server.models.partitioning_criteria_any_of import PartitioningCriteriaAnyOf
from openapi_server.models.pdn_connectivity_stat_report import PdnConnectivityStatReport
from openapi_server.models.pdn_connectivity_status import PdnConnectivityStatus
from openapi_server.models.pdn_connectivity_status_any_of import PdnConnectivityStatusAnyOf
from openapi_server.models.pdu_session_info import PduSessionInfo
from openapi_server.models.pdu_session_information import PduSessionInformation
from openapi_server.models.pdu_session_status import PduSessionStatus
from openapi_server.models.pdu_session_status_any_of import PduSessionStatusAnyOf
from openapi_server.models.pdu_session_status_cfg import PduSessionStatusCfg
from openapi_server.models.pdu_session_type import PduSessionType
from openapi_server.models.pdu_session_type_any_of import PduSessionTypeAnyOf
from openapi_server.models.per_ue_attribute import PerUeAttribute
from openapi_server.models.perf_data import PerfData
from openapi_server.models.performance_data import PerformanceData
from openapi_server.models.performance_data_collection import PerformanceDataCollection
from openapi_server.models.performance_data_info import PerformanceDataInfo
from openapi_server.models.plmn_id import PlmnId
from openapi_server.models.plmn_id_nid import PlmnIdNid
from openapi_server.models.point import Point
from openapi_server.models.point_all_of import PointAllOf
from openapi_server.models.point_altitude import PointAltitude
from openapi_server.models.point_altitude_all_of import PointAltitudeAllOf
from openapi_server.models.point_altitude_uncertainty import PointAltitudeUncertainty
from openapi_server.models.point_altitude_uncertainty_all_of import PointAltitudeUncertaintyAllOf
from openapi_server.models.point_uncertainty_circle import PointUncertaintyCircle
from openapi_server.models.point_uncertainty_circle_all_of import PointUncertaintyCircleAllOf
from openapi_server.models.point_uncertainty_ellipse import PointUncertaintyEllipse
from openapi_server.models.point_uncertainty_ellipse_all_of import PointUncertaintyEllipseAllOf
from openapi_server.models.polygon import Polygon
from openapi_server.models.polygon_all_of import PolygonAllOf
from openapi_server.models.presence_info import PresenceInfo
from openapi_server.models.presence_state import PresenceState
from openapi_server.models.presence_state_any_of import PresenceStateAnyOf
from openapi_server.models.prev_sub_info import PrevSubInfo
from openapi_server.models.problem_details import ProblemDetails
from openapi_server.models.problem_details_analytics_info_request import ProblemDetailsAnalyticsInfoRequest
from openapi_server.models.qoe_metrics_collection import QoeMetricsCollection
from openapi_server.models.qos_requirement import QosRequirement
from openapi_server.models.qos_resource_type import QosResourceType
from openapi_server.models.qos_resource_type_any_of import QosResourceTypeAnyOf
from openapi_server.models.qos_sustainability_info import QosSustainabilityInfo
from openapi_server.models.ranking_criterion import RankingCriterion
from openapi_server.models.rat_freq_information import RatFreqInformation
from openapi_server.models.rat_type import RatType
from openapi_server.models.rat_type_any_of import RatTypeAnyOf
from openapi_server.models.reachability_filter import ReachabilityFilter
from openapi_server.models.reachability_filter_any_of import ReachabilityFilterAnyOf
from openapi_server.models.reachability_for_data_configuration import ReachabilityForDataConfiguration
from openapi_server.models.reachability_for_data_report_config import ReachabilityForDataReportConfig
from openapi_server.models.reachability_for_data_report_config_any_of import ReachabilityForDataReportConfigAnyOf
from openapi_server.models.reachability_for_sms_configuration import ReachabilityForSmsConfiguration
from openapi_server.models.reachability_for_sms_configuration_any_of import ReachabilityForSmsConfigurationAnyOf
from openapi_server.models.reachability_for_sms_report import ReachabilityForSmsReport
from openapi_server.models.reachability_report import ReachabilityReport
from openapi_server.models.red_trans_exp_ordering_criterion import RedTransExpOrderingCriterion
from openapi_server.models.red_trans_exp_ordering_criterion_any_of import RedTransExpOrderingCriterionAnyOf
from openapi_server.models.redundant_transmission_exp_info import RedundantTransmissionExpInfo
from openapi_server.models.redundant_transmission_exp_per_ts import RedundantTransmissionExpPerTS
from openapi_server.models.redundant_transmission_exp_req import RedundantTransmissionExpReq
from openapi_server.models.relative_cartesian_location import RelativeCartesianLocation
from openapi_server.models.report import Report
from openapi_server.models.reporting_information import ReportingInformation
from openapi_server.models.reporting_options import ReportingOptions
from openapi_server.models.requested_context import RequestedContext
from openapi_server.models.resource_usage import ResourceUsage
from openapi_server.models.retainability_threshold import RetainabilityThreshold
from openapi_server.models.rm_info import RmInfo
from openapi_server.models.rm_state import RmState
from openapi_server.models.rm_state_any_of import RmStateAnyOf
from openapi_server.models.roaming_status_report import RoamingStatusReport
from openapi_server.models.route_information import RouteInformation
from openapi_server.models.route_to_location import RouteToLocation
from openapi_server.models.routing_area_id import RoutingAreaId
from openapi_server.models.scheduled_communication_time import ScheduledCommunicationTime
from openapi_server.models.scheduled_communication_time1 import ScheduledCommunicationTime1
from openapi_server.models.scheduled_communication_type import ScheduledCommunicationType
from openapi_server.models.scheduled_communication_type_any_of import ScheduledCommunicationTypeAnyOf
from openapi_server.models.sd_range import SdRange
from openapi_server.models.service_area_id import ServiceAreaId
from openapi_server.models.service_experience_info import ServiceExperienceInfo
from openapi_server.models.service_experience_info1 import ServiceExperienceInfo1
from openapi_server.models.service_experience_info_per_app import ServiceExperienceInfoPerApp
from openapi_server.models.service_experience_info_per_flow import ServiceExperienceInfoPerFlow
from openapi_server.models.service_experience_type import ServiceExperienceType
from openapi_server.models.service_experience_type_any_of import ServiceExperienceTypeAnyOf
from openapi_server.models.service_name import ServiceName
from openapi_server.models.service_name_any_of import ServiceNameAnyOf
from openapi_server.models.sess_inact_timer_for_ue_comm import SessInactTimerForUeComm
from openapi_server.models.slice_load_level_information import SliceLoadLevelInformation
from openapi_server.models.sm_nas_from_smf import SmNasFromSmf
from openapi_server.models.sm_nas_from_ue import SmNasFromUe
from openapi_server.models.smcce_info import SmcceInfo
from openapi_server.models.smcce_info1 import SmcceInfo1
from openapi_server.models.smcce_ue_list import SmcceUeList
from openapi_server.models.smcce_ue_list1 import SmcceUeList1
from openapi_server.models.smf_event import SmfEvent
from openapi_server.models.smf_event_any_of import SmfEventAnyOf
from openapi_server.models.snssai import Snssai
from openapi_server.models.snssai_extension import SnssaiExtension
from openapi_server.models.snssai_tai_mapping import SnssaiTaiMapping
from openapi_server.models.specific_analytics_subscription import SpecificAnalyticsSubscription
from openapi_server.models.specific_data_subscription import SpecificDataSubscription
from openapi_server.models.stationary_indication import StationaryIndication
from openapi_server.models.stationary_indication_any_of import StationaryIndicationAnyOf
from openapi_server.models.supported_gad_shapes import SupportedGADShapes
from openapi_server.models.supported_gad_shapes_any_of import SupportedGADShapesAnyOf
from openapi_server.models.supported_snssai import SupportedSnssai
from openapi_server.models.svc_experience import SvcExperience
from openapi_server.models.tac_range import TacRange
from openapi_server.models.tai import Tai
from openapi_server.models.tai_range import TaiRange
from openapi_server.models.target_area import TargetArea
from openapi_server.models.target_ue_identification import TargetUeIdentification
from openapi_server.models.target_ue_information import TargetUeInformation
from openapi_server.models.threshold_level import ThresholdLevel
from openapi_server.models.time_unit import TimeUnit
from openapi_server.models.time_unit_any_of import TimeUnitAnyOf
from openapi_server.models.time_window import TimeWindow
from openapi_server.models.tnap_id import TnapId
from openapi_server.models.top_application import TopApplication
from openapi_server.models.traffic_characterization import TrafficCharacterization
from openapi_server.models.traffic_descriptor import TrafficDescriptor
from openapi_server.models.traffic_information import TrafficInformation
from openapi_server.models.traffic_profile import TrafficProfile
from openapi_server.models.traffic_profile_any_of import TrafficProfileAnyOf
from openapi_server.models.transaction_info import TransactionInfo
from openapi_server.models.transaction_metric import TransactionMetric
from openapi_server.models.transaction_metric_any_of import TransactionMetricAnyOf
from openapi_server.models.transport_protocol import TransportProtocol
from openapi_server.models.transport_protocol_any_of import TransportProtocolAnyOf
from openapi_server.models.twap_id import TwapId
from openapi_server.models.ueid_ext import UEIdExt
from openapi_server.models.ue_access_behavior_report_item import UeAccessBehaviorReportItem
from openapi_server.models.ue_analytics_context_descriptor import UeAnalyticsContextDescriptor
from openapi_server.models.ue_communication import UeCommunication
from openapi_server.models.ue_communication_collection import UeCommunicationCollection
from openapi_server.models.ue_communication_info import UeCommunicationInfo
from openapi_server.models.ue_in_area_filter import UeInAreaFilter
from openapi_server.models.ue_location_trends_report_item import UeLocationTrendsReportItem
from openapi_server.models.ue_mobility import UeMobility
from openapi_server.models.ue_mobility_collection import UeMobilityCollection
from openapi_server.models.ue_mobility_info import UeMobilityInfo
from openapi_server.models.ue_reachability import UeReachability
from openapi_server.models.ue_reachability_any_of import UeReachabilityAnyOf
from openapi_server.models.ue_trajectory_collection import UeTrajectoryCollection
from openapi_server.models.ue_trajectory_info import UeTrajectoryInfo
from openapi_server.models.ue_type import UeType
from openapi_server.models.ue_type_any_of import UeTypeAnyOf
from openapi_server.models.umt_time import UmtTime
from openapi_server.models.uncertainty_ellipse import UncertaintyEllipse
from openapi_server.models.uncertainty_ellipsoid import UncertaintyEllipsoid
from openapi_server.models.upf_information import UpfInformation
from openapi_server.models.usage_threshold import UsageThreshold
from openapi_server.models.user_data_congestion_collection import UserDataCongestionCollection
from openapi_server.models.user_data_congestion_info import UserDataCongestionInfo
from openapi_server.models.user_location import UserLocation
from openapi_server.models.utra_location import UtraLocation
from openapi_server.models.wlan_ordering_criterion import WlanOrderingCriterion
from openapi_server.models.wlan_ordering_criterion_any_of import WlanOrderingCriterionAnyOf
from openapi_server.models.wlan_per_ss_id_performance_info import WlanPerSsIdPerformanceInfo
from openapi_server.models.wlan_per_ts_performance_info import WlanPerTsPerformanceInfo
from openapi_server.models.wlan_performance_info import WlanPerformanceInfo
from openapi_server.models.wlan_performance_req import WlanPerformanceReq