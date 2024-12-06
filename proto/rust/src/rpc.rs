// This file is @generated by prost-build.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct Empty {}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Handshake {
    #[prost(string, tag = "1")]
    pub client_uuid: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Message {
    #[prost(string, tag = "1")]
    pub event: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub data: ::prost::alloc::string::String,
    #[prost(string, tag = "3")]
    pub message_uuid: ::prost::alloc::string::String,
    #[prost(string, tag = "4")]
    pub client_uuid: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Subscriber {
    #[prost(string, tag = "1")]
    pub event: ::prost::alloc::string::String,
    #[prost(string, tag = "2")]
    pub client_uuid: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Ack {
    #[prost(string, tag = "1")]
    pub data: ::prost::alloc::string::String,
    #[prost(enumeration = "AckCode", tag = "2")]
    pub code: i32,
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum AckCode {
    /// Message successfully received
    Ok = 0,
    /// No matching subscriber for the event
    NoListener = 1,
    /// Queue is full, unable to receive message
    QueueFull = 2,
    /// General server-side error
    InternalError = 3,
    UnknownEvent = 4,
}
impl AckCode {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Ok => "OK",
            Self::NoListener => "NO_LISTENER",
            Self::QueueFull => "QUEUE_FULL",
            Self::InternalError => "INTERNAL_ERROR",
            Self::UnknownEvent => "UNKNOWN_EVENT",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "OK" => Some(Self::Ok),
            "NO_LISTENER" => Some(Self::NoListener),
            "QUEUE_FULL" => Some(Self::QueueFull),
            "INTERNAL_ERROR" => Some(Self::InternalError),
            "UNKNOWN_EVENT" => Some(Self::UnknownEvent),
            _ => None,
        }
    }
}
/// Generated client implementations.
pub mod broker_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct BrokerClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl BrokerClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> BrokerClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> BrokerClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            BrokerClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        /// Server side: handle message reception and dispatch
        ///
        /// Client side: put a message to the queue
        pub async fn enqueue(
            &mut self,
            request: impl tonic::IntoRequest<super::Message>,
        ) -> std::result::Result<tonic::Response<super::Ack>, tonic::Status> {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/rpc.Broker/Enqueue");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("rpc.Broker", "Enqueue"));
            self.inner.unary(req, path, codec).await
        }
        /// Server side: handle subscribing clients
        ///
        /// Client side: subscribe to an event
        pub async fn subscription(
            &mut self,
            request: impl tonic::IntoRequest<super::Subscriber>,
        ) -> std::result::Result<
            tonic::Response<tonic::codec::Streaming<super::Message>>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/rpc.Broker/Subscription");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("rpc.Broker", "Subscription"));
            self.inner.server_streaming(req, path, codec).await
        }
    }
}