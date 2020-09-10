// Objective-C API for talking to github.com/tladesignz/IPtProxy.git Go package.
//   gobind -lang=objc github.com/tladesignz/IPtProxy.git
//
// File is generated by gobind. Do not edit.

#ifndef __IPtProxy_H__
#define __IPtProxy_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


FOUNDATION_EXPORT const int64_t IPtProxyMeekSocksPort;
FOUNDATION_EXPORT const int64_t IPtProxyObfs2SocksPort;
FOUNDATION_EXPORT const int64_t IPtProxyObfs3SocksPort;
FOUNDATION_EXPORT const int64_t IPtProxyObfs4SocksPort;
FOUNDATION_EXPORT const int64_t IPtProxyScramblesuitSocksPort;
FOUNDATION_EXPORT const int64_t IPtProxySnowflakeSocksPort;

/**
 * Start the Obfs4Proxy.
 */
FOUNDATION_EXPORT void IPtProxyStartObfs4Proxy(void);

/**
 * Start the Snowflake client.
 */
FOUNDATION_EXPORT void IPtProxyStartSnowflake(NSString* _Nullable ice, NSString* _Nullable url, NSString* _Nullable front, NSString* _Nullable logFile, BOOL logToStateDir, BOOL keepLocalAddresses, BOOL unsafeLogging, long maxPeers);

#endif
