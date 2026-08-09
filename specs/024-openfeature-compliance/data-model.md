# Data Model: OpenFeature API Compliance

This feature does not introduce any new entities, database tables, or changes to the FlagManagment caching schema. 

The implementation purely focuses on adapting the existing in-memory JSON flag representations (stored in `Client.flags` across the SDKs) to conform to the CNCF OpenFeature `ResolutionDetails` response contract.
