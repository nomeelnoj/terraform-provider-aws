// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inspector2

// Exports for use in tests only.
var (
	ResourceCodeScanningConfiguration = resourceCodeScanningConfiguration
	ResourceDelegatedAdminAccount     = resourceDelegatedAdminAccount
	ResourceFilter                    = newFilterResource
	ResourceMemberAssociation         = resourceMemberAssociation
	ResourceOrganizationConfiguration = resourceOrganizationConfiguration

	FindCodeScanningConfigurationByARN = findCodeScanningConfigurationByARN
	FindDelegatedAdminAccountByID      = findDelegatedAdminAccountByID
	FindFilterByARN                    = findFilterByARN
	FindMemberByAccountID              = findMemberByAccountID
	FindOrganizationConfiguration      = findOrganizationConfiguration

	EnablerID      = enablerID
	ParseEnablerID = parseEnablerID
)
