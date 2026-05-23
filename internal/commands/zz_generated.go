package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ploicloud/cli/internal/client"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ context.Context

func registerGenerated(root *cobra.Command, c *client.Client) {
	gApplications := &cobra.Command{Use: "applications", Short: "Manage applications"}
	root.AddCommand(gApplications)
	cApplicationsAccessibleServices := buildCmd(c, opSpec{
		ID:       "applications.accessible-services",
		Method:   "GET",
		PathTmpl: "/applications/{application}/accessible-services",
		Use:      "accessible-services",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsAccessibleServices)

	cApplicationsAutoscalingHistory := buildCmd(c, opSpec{
		ID:       "applications.autoscaling-history",
		Method:   "GET",
		PathTmpl: "/applications/{application}/autoscaling-history",
		Use:      "autoscaling-history",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{
			{Name: "range", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplications.AddCommand(cApplicationsAutoscalingHistory)

	cApplicationsAutoscalingMetrics := buildCmd(c, opSpec{
		ID:       "applications.autoscaling-metrics",
		Method:   "GET",
		PathTmpl: "/applications/{application}/autoscaling-metrics",
		Use:      "autoscaling-metrics",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsAutoscalingMetrics)

	gApplicationsBasicAuth := &cobra.Command{Use: "basic-auth", Short: "Manage basic-auth"}
	gApplications.AddCommand(gApplicationsBasicAuth)
	cApplicationsBasicAuthShow := buildCmd(c, opSpec{
		ID:       "applications.basic-auth.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/basic-auth",
		Use:      "get",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsBasicAuth.AddCommand(cApplicationsBasicAuthShow)

	cApplicationsBasicAuthUpdate := buildCmd(c, opSpec{
		ID:       "applications.basic-auth.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/basic-auth",
		Use:      "update",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsBasicAuth.AddCommand(cApplicationsBasicAuthUpdate)

	gApplicationsBuildConfig := &cobra.Command{Use: "build-config", Short: "Manage build-config"}
	gApplications.AddCommand(gApplicationsBuildConfig)
	cApplicationsBuildConfigIndex := buildCmd(c, opSpec{
		ID:       "applications.build-config.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/build-config",
		Use:      "list",
		Short:    "Get build configuration",
		Long:     "Retrieve the current build and initialization commands for the application. Build commands\nrun during image construction, while init commands run each time the application starts.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsBuildConfig.AddCommand(cApplicationsBuildConfigIndex)

	cApplicationsBuildConfigUpdate := buildCmd(c, opSpec{
		ID:       "applications.build-config.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/build-config",
		Use:      "update",
		Short:    "Update build configuration",
		Long:     "Update the build and initialization commands for the application. Build commands execute\nduring image build with root permissions. Init commands run with application user\npermissions (UID/GID 33) before starting the main process. Commands are executed sequentially\nand stop on first failure. Updating commands marks the application as needing deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "build_commands", Type: "array:string", Required: false},
			{Name: "init_commands", Type: "array:string", Required: false},
		},
	})
	gApplicationsBuildConfig.AddCommand(cApplicationsBuildConfigUpdate)

	cApplicationsCommandHistory := buildCmd(c, opSpec{
		ID:       "applications.command-history",
		Method:   "GET",
		PathTmpl: "/applications/{application}/command-history",
		Use:      "command-history",
		Short:    "Command history",
		Long:     "Retrieve paginated command execution history for an application. Shows all commands\nthat have been executed, including their status, output, and who executed them.\nResults are ordered by execution time (newest first).",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{
			{Name: "per_page", Type: "integer", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplications.AddCommand(cApplicationsCommandHistory)

	gApplicationsDebug := &cobra.Command{Use: "debug", Short: "Manage debug"}
	gApplications.AddCommand(gApplicationsDebug)
	cApplicationsDebugShow := buildCmd(c, opSpec{
		ID:       "applications.debug.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/debug",
		Use:      "get",
		Short:    "Get debug information",
		Long:     "Retrieve comprehensive debug information for an application including build issues,\ninstance issues, deployment problems, and recommendations.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDebug.AddCommand(cApplicationsDebugShow)

	cApplicationsDebugSummary := buildCmd(c, opSpec{
		ID:       "applications.debug.summary",
		Method:   "GET",
		PathTmpl: "/applications/{application}/debug/summary",
		Use:      "summary",
		Short:    "Get debug summary",
		Long:     "Get a quick summary of the application's health status with critical issue count.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDebug.AddCommand(cApplicationsDebugSummary)

	cApplicationsDeploy := buildCmd(c, opSpec{
		ID:       "applications.deploy",
		Method:   "POST",
		PathTmpl: "/applications/{application}/deploy",
		Use:      "deploy",
		Short:    "Deploy application",
		Long:     "Initiate a new deployment for the application. This will trigger a build process that clones the\nrepository, builds an image, and deploys it to the platform. Any pending or running deployments\nwill be automatically cancelled. The deployment runs asynchronously and progress can be tracked\nvia the deployments endpoint.\n\nWhen skip_build is true, the deployment will reuse the container image from the most recent\nsuccessful deployment instead of building a new one. This is useful for quick redeployments\nwhen only configuration changes are needed. Returns 422 if no previous successful deployment exists.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "branch", Type: "string", Required: false},
			{Name: "skip_build", Type: "boolean", Required: false},
		},
	})
	gApplications.AddCommand(cApplicationsDeploy)

	gApplicationsDeployments := &cobra.Command{Use: "deployments", Short: "Manage deployments"}
	gApplications.AddCommand(gApplicationsDeployments)
	cApplicationsDeploymentsIndex := buildCmd(c, opSpec{
		ID:       "applications.deployments.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/deployments",
		Use:      "list",
		Short:    "List deployments",
		Long:     "Retrieve a paginated list of deployments for an application. Deployments represent the\nbuild and deployment history, showing status, duration, commit information, and who\ntriggered each deployment. Results are ordered by creation date (newest first).",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{
			{Name: "per_page", Type: "integer", Required: false, Desc: ""},
			{Name: "status", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplicationsDeployments.AddCommand(cApplicationsDeploymentsIndex)

	cApplicationsDeploymentsLogs := buildCmd(c, opSpec{
		ID:       "applications.deployments.logs",
		Method:   "GET",
		PathTmpl: "/applications/{application}/deployments/{deployment}/logs",
		Use:      "logs",
		Short:    "Get deployment logs",
		Long:     "Retrieve build and deployment logs for a specific deployment. Includes output from\ngit clone, build commands, image build, and deployment operations.\nLogs are preserved even after the build process is cleaned up.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "deployment", Type: "integer", Required: true, Desc: "The deployment ID"},
		},
		QueryParams: []paramDef{
			{Name: "tail", Type: "integer", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplicationsDeployments.AddCommand(cApplicationsDeploymentsLogs)

	cApplicationsDeploymentsShow := buildCmd(c, opSpec{
		ID:       "applications.deployments.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/deployments/{deployment}",
		Use:      "get",
		Short:    "Get deployment details",
		Long:     "Retrieve detailed information about a specific deployment including full commit details,\nexecution timeline, and error information if the deployment failed. Use this endpoint\nto check deployment status and troubleshoot failures.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "deployment", Type: "integer", Required: true, Desc: "The deployment ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDeployments.AddCommand(cApplicationsDeploymentsShow)

	cApplicationsDestroy := buildCmd(c, opSpec{
		ID:       "applications.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}",
		Use:      "delete",
		Short:    "Delete application",
		Long:     "Permanently delete an application and all associated resources including services, deployments,\nsecrets, and platform resources. This action cannot be undone. The application environment\nwill be removed from the platform.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsDestroy)

	gApplicationsDevSession := &cobra.Command{Use: "dev-session", Short: "Manage dev-session"}
	gApplications.AddCommand(gApplicationsDevSession)
	cApplicationsDevSessionDestroy := buildCmd(c, opSpec{
		ID:       "applications.dev-session.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/dev-session",
		Use:      "delete",
		Short:    "Destroy dev session",
		Long:     "Destroy the active dev session for an application. This removes the dev container\nand scales the production deployment back up.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDevSession.AddCommand(cApplicationsDevSessionDestroy)

	cApplicationsDevSessionShow := buildCmd(c, opSpec{
		ID:       "applications.dev-session.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/dev-session",
		Use:      "get",
		Short:    "Show dev session",
		Long:     "Get the active dev session for an application, if one exists.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDevSession.AddCommand(cApplicationsDevSessionShow)

	cApplicationsDevSessionStore := buildCmd(c, opSpec{
		ID:       "applications.dev-session.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/dev-session",
		Use:      "create",
		Short:    "Create dev session",
		Long:     "Create a new dev session for an application. This scales down the production deployment\nand creates a dev container with an AI agent for live editing.\nOnly one dev session per application is allowed.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDevSession.AddCommand(cApplicationsDevSessionStore)

	cApplicationsDismissNotification := buildCmd(c, opSpec{
		ID:       "applications.dismiss-notification",
		Method:   "POST",
		PathTmpl: "/applications/{application}/dismiss-notification",
		Use:      "dismiss-notification",
		Short:    "Dismiss notifications",
		Long:     "Dismiss getting started or deployment needed notifications.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "type", Type: "string", Required: true},
		},
	})
	gApplications.AddCommand(cApplicationsDismissNotification)

	gApplicationsDomains := &cobra.Command{Use: "domains", Short: "Manage domains"}
	gApplications.AddCommand(gApplicationsDomains)
	cApplicationsDomainsDestroy := buildCmd(c, opSpec{
		ID:       "applications.domains.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/domains/{domain}",
		Use:      "delete",
		Short:    "Remove application domain",
		Long:     "Remove a custom domain from the application. By default the application must be\nredeployed for the change to take effect. Set `force_apply` to `true` to immediately\nupdate routing — if this was the last custom domain, the application's preview domain\nis automatically restored.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "domain", Type: "integer", Required: true, Desc: "The domain ID"},
		},
		QueryParams: []paramDef{
			{Name: "force_apply", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplicationsDomains.AddCommand(cApplicationsDomainsDestroy)

	cApplicationsDomainsIndex := buildCmd(c, opSpec{
		ID:       "applications.domains.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/domains",
		Use:      "list",
		Short:    "List application domains",
		Long:     "Retrieve all custom domains associated with an application. Does not include the primary\napplication domain (e.g., app-name.test.ploi.it). Each domain automatically receives an\nSSL certificate via Let's Encrypt after deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDomains.AddCommand(cApplicationsDomainsIndex)

	cApplicationsDomainsSslStatus := buildCmd(c, opSpec{
		ID:       "applications.domains.ssl-status",
		Method:   "GET",
		PathTmpl: "/applications/{application}/domains/{domain}/ssl-status",
		Use:      "ssl-status",
		Short:    "Get domain SSL status",
		Long:     "Retrieve the SSL certificate status for a specific domain. Shows certificate issuer,\nexpiration date, and whether auto-renewal is enabled. Certificates are managed by\ncert-manager and automatically renewed 30 days before expiration.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "domain", Type: "integer", Required: true, Desc: "The domain ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsDomains.AddCommand(cApplicationsDomainsSslStatus)

	cApplicationsDomainsStore := buildCmd(c, opSpec{
		ID:       "applications.domains.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/domains",
		Use:      "create",
		Short:    "Add application domain",
		Long:     "Add a custom domain to the application. The domain must be properly configured with DNS\npointing to the platform's load balancer before adding. SSL certificates are automatically\nprovisioned via Let's Encrypt during the next deployment, unless `force_apply` is set\nto `true` — in which case the domain is activated and SSL provisioning begins immediately\nwithout requiring a deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "domain", Type: "string", Required: true},
			{Name: "force_apply", Type: "boolean", Required: false},
			{Name: "skip_dns_check", Type: "boolean", Required: false},
		},
	})
	gApplicationsDomains.AddCommand(cApplicationsDomainsStore)

	cApplicationsExec := buildCmd(c, opSpec{
		ID:       "applications.exec",
		Method:   "POST",
		PathTmpl: "/applications/{application}/exec",
		Use:      "exec",
		Short:    "Execute command",
		Long:     "Execute a command in a running instance. The command is executed via the container's\nshell and the output is returned. All command executions are logged for audit purposes.\nCommands have a 120-second timeout.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "command", Type: "string", Required: true},
			{Name: "instance", Type: "string", Required: true},
		},
	})
	gApplications.AddCommand(cApplicationsExec)

	cApplicationsExportYaml := buildCmd(c, opSpec{
		ID:       "applications.export-yaml",
		Method:   "GET",
		PathTmpl: "/applications/{application}/export-yaml",
		Use:      "export-yaml",
		Short:    "Export application YAML definition",
		Long:     "Export the current application configuration as a YAML file compatible with the\ninfrastructure apply endpoint. The exported YAML includes application settings,\nrepository configuration, runtime, commands, services, secrets, domains, and volumes.\n\nThis YAML can be re-imported via `POST /api/v1/infrastructure/apply` to recreate\nor update the application configuration.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsExportYaml)

	cApplicationsIndex := buildCmd(c, opSpec{
		ID:         "applications.index",
		Method:     "GET",
		PathTmpl:   "/applications",
		Use:        "list",
		Short:      "List applications",
		Long:       "Retrieve a paginated list of applications for the current team. Applications can be filtered by status\nand searched by name. The response includes basic application information along with associated domains\nand services.",
		PathParams: []paramDef{},
		QueryParams: []paramDef{
			{Name: "per_page", Type: "integer", Required: false, Desc: ""},
			{Name: "search", Type: "string", Required: false, Desc: ""},
			{Name: "status", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplications.AddCommand(cApplicationsIndex)

	cApplicationsInstances := buildCmd(c, opSpec{
		ID:       "applications.instances",
		Method:   "GET",
		PathTmpl: "/applications/{application}/instances",
		Use:      "instances",
		Short:    "List running instances",
		Long:     "Retrieve all running pod instances for an application. Returns both web instances\nand scheduler instances (if scheduler is enabled). Each instance can be used\nas a target for command execution.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsInstances)

	cApplicationsLogs := buildCmd(c, opSpec{
		ID:       "applications.logs",
		Method:   "GET",
		PathTmpl: "/applications/{application}/logs",
		Use:      "logs",
		Short:    "Get application logs",
		Long:     "Retrieve logs from the main application instance. Logs are fetched directly from the\napplication instance and include stdout/stderr output from the application. Use the tail\nparameter to limit the amount of data returned for large log files.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{
			{Name: "tail", Type: "integer", Required: false, Desc: ""},
			{Name: "since", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplications.AddCommand(cApplicationsLogs)

	gApplicationsNetworks := &cobra.Command{Use: "networks", Short: "Manage networks"}
	gApplications.AddCommand(gApplicationsNetworks)
	cApplicationsNetworksIndex := buildCmd(c, opSpec{
		ID:       "applications.networks.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/networks",
		Use:      "list",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsNetworks.AddCommand(cApplicationsNetworksIndex)

	cApplicationsNetworksUpdate := buildCmd(c, opSpec{
		ID:       "applications.networks.update",
		Method:   "PUT",
		PathTmpl: "/applications/{application}/networks",
		Use:      "update",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "network_ids", Type: "array:integer", Required: false},
		},
	})
	gApplicationsNetworks.AddCommand(cApplicationsNetworksUpdate)

	gApplicationsPhpConfig := &cobra.Command{Use: "php-config", Short: "Manage php-config"}
	gApplications.AddCommand(gApplicationsPhpConfig)
	cApplicationsPhpConfigUpdate := buildCmd(c, opSpec{
		ID:       "applications.php-config.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/php-config",
		Use:      "update",
		Short:    "Update PHP configuration",
		Long:     "Configure PHP version, extensions, and runtime settings for the application. Changes require\na new deployment to take effect. PHP settings are applied via custom php.ini configuration.\nExtensions are installed using the install-php-extensions tool during build.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "php_extensions", Type: "array:string", Required: false},
			{Name: "php_image_variant", Type: "string", Required: false},
			{Name: "php_settings", Type: "array:string", Required: false},
			{Name: "php_version", Type: "string", Required: false},
			{Name: "start_command", Type: "string", Required: false},
		},
	})
	gApplicationsPhpConfig.AddCommand(cApplicationsPhpConfigUpdate)

	gApplicationsRepository := &cobra.Command{Use: "repository", Short: "Manage repository"}
	gApplications.AddCommand(gApplicationsRepository)
	cApplicationsRepositoryUpdate := buildCmd(c, opSpec{
		ID:       "applications.repository.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/repository",
		Use:      "update",
		Short:    "Update repository settings",
		Long:     "Update the repository configuration for an application.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "default_branch", Type: "string", Required: false},
			{Name: "repository_name", Type: "string", Required: false},
			{Name: "repository_owner", Type: "string", Required: false},
			{Name: "repository_url", Type: "string", Required: false},
			{Name: "social_account_id", Type: "integer", Required: false},
		},
	})
	gApplicationsRepository.AddCommand(cApplicationsRepositoryUpdate)

	gApplicationsResources := &cobra.Command{Use: "resources", Short: "Manage resources"}
	gApplications.AddCommand(gApplicationsResources)
	cApplicationsResourcesUpdate := buildCmd(c, opSpec{
		ID:       "applications.resources.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/resources",
		Use:      "update",
		Short:    "Update resource allocation",
		Long:     "Update resource allocation settings for an application. This endpoint requires a RAM value\nand automatically calculates CPU allocation based on the predefined RAM_TO_CPU_MAP.\nChanges require a new deployment to take effect.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "memory_request", Type: "string", Required: true},
			{Name: "replicas", Type: "integer", Required: false},
		},
	})
	gApplicationsResources.AddCommand(cApplicationsResourcesUpdate)

	cApplicationsResume := buildCmd(c, opSpec{
		ID:       "applications.resume",
		Method:   "POST",
		PathTmpl: "/applications/{application}/resume",
		Use:      "resume",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsResume)

	cApplicationsRollback := buildCmd(c, opSpec{
		ID:       "applications.rollback",
		Method:   "POST",
		PathTmpl: "/applications/{application}/rollback",
		Use:      "rollback",
		Short:    "Rollback application",
		Long:     "Rollback an application to a previous successful deployment by rebuilding from its commit SHA.\nThis triggers a full rebuild targeting the specific commit rather than reusing an old image.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "deployment_id", Type: "integer", Required: true},
		},
	})
	gApplications.AddCommand(cApplicationsRollback)

	gApplicationsSecrets := &cobra.Command{Use: "secrets", Short: "Manage secrets"}
	gApplications.AddCommand(gApplicationsSecrets)
	cApplicationsSecretsDestroy := buildCmd(c, opSpec{
		ID:       "applications.secrets.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/secrets/{key}",
		Use:      "delete",
		Short:    "Delete application secret",
		Long:     "Remove an environment variable from the application. This action cannot be undone.\nDeleting a secret automatically marks the application as needing deployment. The\nsecret will be removed from the application environment on the next deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "key", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSecrets.AddCommand(cApplicationsSecretsDestroy)

	cApplicationsSecretsIndex := buildCmd(c, opSpec{
		ID:       "applications.secrets.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/secrets",
		Use:      "list",
		Short:    "List application secrets",
		Long:     "Retrieve all environment variables configured for an application. Secret values are always\nmasked as '********' in API responses for security. Secrets are encrypted at rest in the\ndatabase. Changes to secrets require a new deployment to take effect.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSecrets.AddCommand(cApplicationsSecretsIndex)

	cApplicationsSecretsInjected := buildCmd(c, opSpec{
		ID:       "applications.secrets.injected",
		Method:   "GET",
		PathTmpl: "/applications/{application}/secrets/injected",
		Use:      "injected",
		Short:    "Get injected secrets",
		Long:     "Retrieve environment variables automatically injected by application services such as\ndatabases, caches, and other services. These are not user-defined secrets but are\ngenerated by the platform based on the services attached to the application. Values\nare shown in plain text as they are not user secrets.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSecrets.AddCommand(cApplicationsSecretsInjected)

	cApplicationsSecretsStore := buildCmd(c, opSpec{
		ID:       "applications.secrets.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/secrets",
		Use:      "create",
		Short:    "Create application secret",
		Long:     "Add a new environment variable to the application. Secret keys must be uppercase letters,\nnumbers, and underscores only. Values are encrypted before storage. All secrets are\npassed as build arguments during image construction and injected as environment variables\nat runtime.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "key", Type: "string", Required: true},
			{Name: "value", Type: "string", Required: true},
		},
	})
	gApplicationsSecrets.AddCommand(cApplicationsSecretsStore)

	cApplicationsSecretsUpdate := buildCmd(c, opSpec{
		ID:       "applications.secrets.update",
		Method:   "PUT",
		PathTmpl: "/applications/{application}/secrets/{key}",
		Use:      "update",
		Short:    "Update application secret",
		Long:     "Update the value of an existing environment variable. The key cannot be changed - to rename\na secret, delete and recreate it. Updating a secret automatically marks the application as\nneeding deployment. The new value will be encrypted before storage.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "key", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "value", Type: "string", Required: true},
		},
	})
	gApplicationsSecrets.AddCommand(cApplicationsSecretsUpdate)

	gApplicationsServices := &cobra.Command{Use: "services", Short: "Manage services"}
	gApplications.AddCommand(gApplicationsServices)
	cApplicationsServicesConnectionFormats := buildCmd(c, opSpec{
		ID:       "applications.services.connection-formats",
		Method:   "GET",
		PathTmpl: "/applications/{application}/services/{service}/connection-formats",
		Use:      "connection-formats",
		Short:    "Get connection formats for a service",
		Long:     "Retrieve available connection string formats for a service with debug access enabled.\nSupports multiple formats including standard URIs, JDBC connections, and CLI commands\nfor different database types.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesConnectionFormats)

	gApplicationsServicesDebug := &cobra.Command{Use: "debug", Short: "Manage debug"}
	gApplicationsServices.AddCommand(gApplicationsServicesDebug)
	cApplicationsServicesDebugDisable := buildCmd(c, opSpec{
		ID:       "applications.services.debug.disable",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/services/{service}/debug",
		Use:      "disable",
		Short:    "Disable debug access",
		Long:     "Remove debug access to a service. This deletes the temporary external service and revokes\nexternal access. Any active connections will be terminated.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServicesDebug.AddCommand(cApplicationsServicesDebugDisable)

	cApplicationsServicesDebugEnable := buildCmd(c, opSpec{
		ID:       "applications.services.debug.enable",
		Method:   "POST",
		PathTmpl: "/applications/{application}/services/{service}/debug",
		Use:      "enable",
		Short:    "Enable debug access",
		Long:     "Create temporary external access to a service for debugging purposes. This creates an external\nservice with a random port that allows direct connection from outside the platform. Access\nautomatically expires after the specified duration. Useful for database administration or\ntroubleshooting.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServicesDebug.AddCommand(cApplicationsServicesDebugEnable)

	cApplicationsServicesDestroy := buildCmd(c, opSpec{
		ID:       "applications.services.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/services/{service}",
		Use:      "delete",
		Short:    "Delete application service",
		Long:     "Permanently delete a service and all associated data. This will remove the service,\npersistent volumes, and any stored data. This action cannot be undone. The service will be\nmarked for deletion and removed asynchronously.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesDestroy)

	cApplicationsServicesEnvironment := buildCmd(c, opSpec{
		ID:       "applications.services.environment",
		Method:   "GET",
		PathTmpl: "/applications/{application}/services/{service}/environment",
		Use:      "environment",
		Short:    "Get service environment variables",
		Long:     "Retrieve the environment variables that this service automatically injects into the main\napplication. These variables provide connection details and credentials for\naccessing the service from the application.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesEnvironment)

	cApplicationsServicesIndex := buildCmd(c, opSpec{
		ID:       "applications.services.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/services",
		Use:      "list",
		Short:    "List application services",
		Long:     "Retrieve all services associated with an application. Services include databases (MySQL, PostgreSQL, MongoDB),\ncaches (Redis, Valkey), message queues (RabbitMQ), object storage (MinIO), background workers, and SFTP servers.\nEach service runs in the same application environment and automatically injects environment\nvariables into the main application.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesIndex)

	cApplicationsServicesLogs := buildCmd(c, opSpec{
		ID:       "applications.services.logs",
		Method:   "GET",
		PathTmpl: "/applications/{application}/services/{service}/logs",
		Use:      "logs",
		Short:    "Get service logs",
		Long:     "Retrieve logs from a specific application service (database, cache, worker, etc.).\nIncludes service-specific output such as database queries, cache operations, or\nworker job processing.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{
			{Name: "tail", Type: "integer", Required: false, Desc: ""},
			{Name: "since", Type: "string", Required: false, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesLogs)

	cApplicationsServicesRestart := buildCmd(c, opSpec{
		ID:       "applications.services.restart",
		Method:   "POST",
		PathTmpl: "/applications/{application}/services/{service}/restart",
		Use:      "restart",
		Short:    "Restart service",
		Long:     "Restart a service. This causes a brief downtime\nwhile the service restarts. Useful for applying configuration changes or recovering from issues.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesRestart)

	cApplicationsServicesStore := buildCmd(c, opSpec{
		ID:       "applications.services.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/services",
		Use:      "create",
		Short:    "Create application service",
		Long:     "Add a new service to an application. Services run as separate services in the same environment and\nautomatically inject connection details as environment variables. Database services create persistent\nvolumes for data storage. Worker services share the same application image and volumes as the main application.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "name", Type: "string", Required: true},
			{Name: "settings", Type: "object", Required: false},
			{Name: "type", Type: "string", Required: true},
			{Name: "version", Type: "string", Required: false},
		},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesStore)

	cApplicationsServicesUpdate := buildCmd(c, opSpec{
		ID:       "applications.services.update",
		Method:   "PUT",
		PathTmpl: "/applications/{application}/services/{service}",
		Use:      "update",
		Short:    "Update service configuration",
		Long:     "Update configuration for an existing service. Only certain settings can be modified after creation.\nChanges may require a service restart to take effect. Worker services can update their command,\nwhile other services have limited update capabilities.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "service", Type: "integer", Required: true, Desc: "The service ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "settings", Type: "object", Required: false},
			{Name: "version", Type: "string", Required: false},
		},
	})
	gApplicationsServices.AddCommand(cApplicationsServicesUpdate)

	gApplicationsSettings := &cobra.Command{Use: "settings", Short: "Manage settings"}
	gApplications.AddCommand(gApplicationsSettings)
	cApplicationsSettingsUpdate := buildCmd(c, opSpec{
		ID:       "applications.settings.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/settings",
		Use:      "update",
		Short:    "Update application settings",
		Long:     "Update various application settings including health checks and scheduler configuration.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "health_check_path", Type: "string", Required: false},
			{Name: "nodejs_version", Type: "string", Required: false},
			{Name: "scheduler_enabled", Type: "boolean", Required: false},
			{Name: "webroot_path", Type: "string", Required: false},
		},
	})
	gApplicationsSettings.AddCommand(cApplicationsSettingsUpdate)

	cApplicationsShow := buildCmd(c, opSpec{
		ID:       "applications.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}",
		Use:      "get",
		Short:    "Get application details",
		Long:     "Retrieve comprehensive information about a specific application including its configuration, domains,\nservices, and secrets. Only accessible to team members with appropriate permissions.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsShow)

	gApplicationsSource := &cobra.Command{Use: "source", Short: "Manage source"}
	gApplications.AddCommand(gApplicationsSource)
	cApplicationsSourceDownload := buildCmd(c, opSpec{
		ID:       "applications.source.download",
		Method:   "GET",
		PathTmpl: "/applications/{application}/source/download",
		Use:      "download",
		Short:    "Download source archive",
		Long:     "Download the source archive for an application. This endpoint uses signed URLs\nfor secure access by build pods during the deployment process.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSource.AddCommand(cApplicationsSourceDownload)

	cApplicationsSourceSignedUpload := buildCmd(c, opSpec{
		ID:       "applications.source.signed-upload",
		Method:   "POST",
		PathTmpl: "/applications/{application}/source/signed-upload",
		Use:      "signed-upload",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSource.AddCommand(cApplicationsSourceSignedUpload)

	cApplicationsSourceUpload := buildCmd(c, opSpec{
		ID:       "applications.source.upload",
		Method:   "POST",
		PathTmpl: "/applications/{application}/source",
		Use:      "upload",
		Short:    "Upload source archive",
		Long:     "Upload a source archive (.zip or .tar.gz) for an application that uses upload-based deployments\ninstead of a Git repository. The uploaded archive will be used as the source code for the next deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsSource.AddCommand(cApplicationsSourceUpload)

	cApplicationsStatus := buildCmd(c, opSpec{
		ID:       "applications.status",
		Method:   "GET",
		PathTmpl: "/applications/{application}/status",
		Use:      "status",
		Short:    "Get application status",
		Long:     "Get detailed status information including SSL certificate status and deployment readiness.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsStatus)

	cApplicationsStore := buildCmd(c, opSpec{
		ID:          "applications.store",
		Method:      "POST",
		PathTmpl:    "/applications",
		Use:         "create",
		Short:       "Create application",
		Long:        "Create a new containerized application in the current team. The application will be deployed to the cloud\nplatform with automatic SSL certificate provisioning. Laravel applications receive default build and init\ncommands automatically. WordPress applications can either use a Git repository or a persistent volume for\nfile storage.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "application_type", Type: "string", Required: true},
			{Name: "application_version", Type: "string", Required: false},
			{Name: "default_branch", Type: "string", Required: false},
			{Name: "name", Type: "string", Required: true},
			{Name: "php_version", Type: "string", Required: false},
			{Name: "provider", Type: "string", Required: false},
			{Name: "region", Type: "string", Required: false},
			{Name: "repository_name", Type: "string", Required: false},
			{Name: "repository_owner", Type: "string", Required: false},
			{Name: "repository_url", Type: "string", Required: false},
			{Name: "scheduled_deletion_at", Type: "string", Required: false},
			{Name: "social_account_id", Type: "integer", Required: false},
			{Name: "source_type", Type: "string", Required: false},
			{Name: "statamic_use_git_integration", Type: "boolean", Required: false},
			{Name: "webroot_path", Type: "string", Required: false},
			{Name: "wordpress_use_volume", Type: "boolean", Required: false},
		},
	})
	gApplications.AddCommand(cApplicationsStore)

	cApplicationsSuspend := buildCmd(c, opSpec{
		ID:       "applications.suspend",
		Method:   "POST",
		PathTmpl: "/applications/{application}/suspend",
		Use:      "suspend",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplications.AddCommand(cApplicationsSuspend)

	cApplicationsUpdate := buildCmd(c, opSpec{
		ID:       "applications.update",
		Method:   "PUT",
		PathTmpl: "/applications/{application}",
		Use:      "update",
		Short:    "Update application",
		Long:     "Update application configuration including name, build commands, and initialization commands.\nModifying build or init commands will automatically mark the application as needing deployment.\nChanges to commands will take effect on the next deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "build_commands", Type: "array:string", Required: false},
			{Name: "init_commands", Type: "array:string", Required: false},
			{Name: "scheduled_deletion_at", Type: "string", Required: false},
			{Name: "webroot_path", Type: "string", Required: false},
		},
	})
	gApplications.AddCommand(cApplicationsUpdate)

	gApplicationsVolumes := &cobra.Command{Use: "volumes", Short: "Manage volumes"}
	gApplications.AddCommand(gApplicationsVolumes)
	cApplicationsVolumesDestroy := buildCmd(c, opSpec{
		ID:       "applications.volumes.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/volumes/{volume}",
		Use:      "delete",
		Short:    "Delete a volume",
		Long:     "Remove a custom volume from the application. Only application volumes can be deleted;\nservice volumes are managed through the services API. Deleting a volume marks the\napplication as needing deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "volume", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsVolumes.AddCommand(cApplicationsVolumesDestroy)

	cApplicationsVolumesIndex := buildCmd(c, opSpec{
		ID:       "applications.volumes.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/volumes",
		Use:      "list",
		Short:    "List application volumes",
		Long:     "Retrieve all volumes associated with an application. This includes both service volumes\n(from databases, caches, etc.) and custom application volumes. Service volumes have\ntype: 'service' while custom volumes have type: 'application'.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsVolumes.AddCommand(cApplicationsVolumesIndex)

	cApplicationsVolumesShow := buildCmd(c, opSpec{
		ID:       "applications.volumes.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/volumes/{volume}",
		Use:      "get",
		Short:    "Get volume details",
		Long:     "Retrieve detailed information about a specific volume including size, usage,\nand status. The volumeId parameter can be either the PVC name for service volumes\n(e.g., \"db-pvc\") or the numeric ID for application volumes.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "volume", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsVolumes.AddCommand(cApplicationsVolumesShow)

	cApplicationsVolumesStore := buildCmd(c, opSpec{
		ID:       "applications.volumes.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/volumes",
		Use:      "create",
		Short:    "Create a new volume",
		Long:     "Add a custom volume to the application. The volume will be created as a persistent\nvolume claim and mounted at the specified path. Creating a volume marks the\napplication as needing deployment.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "name", Type: "string", Required: true},
			{Name: "path", Type: "string", Required: true},
			{Name: "size", Type: "integer", Required: true},
		},
	})
	gApplicationsVolumes.AddCommand(cApplicationsVolumesStore)

	cApplicationsVolumesUpdate := buildCmd(c, opSpec{
		ID:       "applications.volumes.update",
		Method:   "PATCH",
		PathTmpl: "/applications/{application}/volumes/{volume}",
		Use:      "update",
		Short:    "Update volume (resize)",
		Long:     "Resize a volume to a larger size. The resize operation is performed asynchronously\nand may cause brief downtime while the service is restarted. The new size must be\nlarger than the current size and cannot exceed 1000 GB.",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "volume", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "size", Type: "integer", Required: true},
		},
	})
	gApplicationsVolumes.AddCommand(cApplicationsVolumesUpdate)

	gApplicationsWildcardDomains := &cobra.Command{Use: "wildcard-domains", Short: "Manage wildcard-domains"}
	gApplications.AddCommand(gApplicationsWildcardDomains)
	cApplicationsWildcardDomainsActivate := buildCmd(c, opSpec{
		ID:       "applications.wildcard-domains.activate",
		Method:   "POST",
		PathTmpl: "/applications/{application}/wildcard-domains/{wildcardDomain}/activate",
		Use:      "activate",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "wildcardDomain", Type: "integer", Required: true, Desc: "The wildcard domain ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsWildcardDomains.AddCommand(cApplicationsWildcardDomainsActivate)

	cApplicationsWildcardDomainsDestroy := buildCmd(c, opSpec{
		ID:       "applications.wildcard-domains.destroy",
		Method:   "DELETE",
		PathTmpl: "/applications/{application}/wildcard-domains/{wildcardDomain}",
		Use:      "delete",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "wildcardDomain", Type: "integer", Required: true, Desc: "The wildcard domain ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsWildcardDomains.AddCommand(cApplicationsWildcardDomainsDestroy)

	cApplicationsWildcardDomainsIndex := buildCmd(c, opSpec{
		ID:       "applications.wildcard-domains.index",
		Method:   "GET",
		PathTmpl: "/applications/{application}/wildcard-domains",
		Use:      "list",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsWildcardDomains.AddCommand(cApplicationsWildcardDomainsIndex)

	cApplicationsWildcardDomainsShow := buildCmd(c, opSpec{
		ID:       "applications.wildcard-domains.show",
		Method:   "GET",
		PathTmpl: "/applications/{application}/wildcard-domains/{wildcardDomain}",
		Use:      "get",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
			{Name: "wildcardDomain", Type: "integer", Required: true, Desc: "The wildcard domain ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gApplicationsWildcardDomains.AddCommand(cApplicationsWildcardDomainsShow)

	cApplicationsWildcardDomainsStore := buildCmd(c, opSpec{
		ID:       "applications.wildcard-domains.store",
		Method:   "POST",
		PathTmpl: "/applications/{application}/wildcard-domains",
		Use:      "create",
		Short:    "",
		PathParams: []paramDef{
			{Name: "application", Type: "integer", Required: true, Desc: "The application ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "domain", Type: "string", Required: true},
		},
	})
	gApplicationsWildcardDomains.AddCommand(cApplicationsWildcardDomainsStore)

	gCluster := &cobra.Command{Use: "cluster", Short: "Manage cluster"}
	root.AddCommand(gCluster)
	cClusterIpRanges := buildCmd(c, opSpec{
		ID:          "cluster.ip-ranges",
		Method:      "GET",
		PathTmpl:    "/cluster/ip-ranges",
		Use:         "ip-ranges",
		Short:       "Get IP ranges for all nodes grouped by region",
		Long:        "Returns aggregated /24 IP ranges for all Kubernetes nodes, grouped by region.\nThis endpoint is useful for clients who need to whitelist Ploi Cloud IP addresses\nin their external services or firewalls.\n\n## Response Format\n\n```json\n{\n  \"regions\": {\n    \"ams1\": [\n      \"185.52.172.0/24\",\n      \"185.52.173.0/24\"\n    ]\n  }\n}\n```",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gCluster.AddCommand(cClusterIpRanges)

	gInfrastructure := &cobra.Command{Use: "infrastructure", Short: "Manage infrastructure"}
	root.AddCommand(gInfrastructure)
	cInfrastructureApply := buildCmd(c, opSpec{
		ID:          "infrastructure.apply",
		Method:      "POST",
		PathTmpl:    "/infrastructure/apply",
		Use:         "apply",
		Short:       "Apply infrastructure configuration from YAML",
		Long:        "This endpoint processes infrastructure-as-code definitions to create or update\napplications and their associated resources on the Ploi Cloud platform.\n\n## YAML Format\n\nThe YAML configuration should follow this structure:\n\n```yaml\napiVersion: v1\nkind: Infrastructure\nmetadata:\n  name: my-app          # Application identifier\n  team: 1               # Team ID that owns this infrastructure\nspec:\n  application:\n    type: laravel       # Application type: laravel, nodejs, wordpress\n    version: \"12\"       # Framework version\n    label: My App       # Display name (optional)\n    tags:               # Tags for organizing apps (optional)\n      - production\n      - api\n    repository:\n      url: https://github.com/user/repo\n      owner: user\n      name: repo\n      branch: main\n    runtime:\n      php_version: 8.4           # PHP version (for PHP apps)\n      nodejs_version: \"24\"       # Node.js version\n    commands:\n      build:                     # Commands run during build\n        - npm ci\n        - npm run build\n      init:                      # Commands run before app starts\n        - php artisan migrate\n      start: npm start           # Override start command\n    settings:\n      health_check_path: /health\n      scheduler_enabled: true    # Enable Laravel scheduler\n      replicas: 3                # Number of replicas\n      memory: 1024Mi             # Memory limit\n      scheduled_deletion_at: \"2026-12-01T00:00:00Z\"  # Optional ISO 8601 timestamp; the application and all of its data will be permanently deleted at this time. Owners receive warning emails 7 and 1 day before. Removing this field clears the schedule.\n    php:                         # PHP-specific settings\n      extensions:\n        - ldap\n        - imagick\n      settings:\n        - memory_limit=512M\n        - max_execution_time=60\n    security:                    # Optional: Nginx security configuration\n      enabled: true              # Enable security headers\n      headers:                   # Custom HTTP security headers (optional when enabled=true)\n                                 # Example OWASP-recommended headers (from UI \"Set default values\" button):\n        Cache-Control: \"no-store, max-age=0\"\n        Content-Security-Policy: \"default-src 'self'; form-action 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; upgrade-insecure-requests\"\n        Cross-Origin-Embedder-Policy: require-corp\n        Cross-Origin-Opener-Policy: same-origin\n        Cross-Origin-Resource-Policy: same-origin\n        Permissions-Policy: \"accelerometer=(),ambient-light-sensor=(),autoplay=(),battery=(),camera=(),display-capture=(),document-domain=(),encrypted-media=(),execution-while-not-rendered=(),execution-while-out-of-viewport=(),fullscreen=(),gamepad=(),geolocation=(),gyroscope=(),hid=(),idle-detection=(),local-fonts=(),magnetometer=(),microphone=(),midi=(),payment=(),picture-in-picture=(),publickey-credentials-get=(),screen-wake-lock=(),serial=(),speaker-selection=(),usb=(),web-share=(),xr-spatial-tracking=()\"\n        Referrer-Policy: no-referrer\n        Strict-Transport-Security: \"max-age=31536000; includeSubDomains\"\n        X-Content-Type-Options: nosniff\n        X-Frame-Options: deny\n        X-Permitted-Cross-Domain-Policies: none\n      ssl_protocols: \"TLSv1.2 TLSv1.3\"    # SSL/TLS protocols (optional)\n      ssl_ciphers: \"ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305\"  # SSL cipher suites (optional)\n\n  domains:\n    - domain: app.example.com\n    - domain: www.example.com\n\n  secrets:                       # Environment variables\n    - key: APP_KEY\n      value: base64:your-app-key\n    - key: DB_PASSWORD\n      value: secret123\n\n  volumes:                       # Persistent volumes\n    - name: storage\n      mount_path: /var/www/html/storage\n      volume_size: 10            # Size in GB\n\n  services:                      # Database/cache services\n    - name: database\n      type: mysql                # mysql, postgresql, mongodb, redis, valkey, rabbitmq, minio, sftp\n      version: \"8.0\"\n      memory: 2Gi\n      volume_size: 20Gi\n      settings:                  # Service-specific settings\n        extensions:\n          - postgis\n\n    - name: cache\n      type: redis\n      version: \"7.2\"\n      memory: 512Mi\n      volume_size: 1Gi\n\n    - name: queue-worker        # Worker service\n      type: worker\n      memory: 1Gi\n      command: php artisan queue:work\n\n  container_services:            # Pre-built container services\n    - name: pdf-generator\n      type: gotenberg            # gotenberg, chrome-headless, clickhouse\n      version: \"8\"\n      memory: 1Gi\n```\n\n## Query Parameters\n\n- `dry_run` (boolean, default: false) - When true, shows what changes would be made without applying them\n- `auto_deploy` (boolean, default: true) - When true, automatically deploys the application after changes\n\n## Example cURL Request\n\n```bash\n# Normal deployment\ncurl -X POST https://api.ploi.cloud/api/v1/infrastructure/apply \\\n  -H \"Authorization: Bearer YOUR_API_TOKEN\" \\\n  -H \"Content-Type: application/yaml\" \\\n  --data-binary @ploi.yaml\n\n# Dry run mode - preview changes without applying\ncurl -X POST \"https://api.ploi.cloud/api/v1/infrastructure/apply?dry_run=true\" \\\n  -H \"Authorization: Bearer YOUR_API_TOKEN\" \\\n  -H \"Content-Type: application/yaml\" \\\n  --data-binary @ploi.yaml\n\n# Apply changes without automatic deployment\ncurl -X POST \"https://api.ploi.cloud/api/v1/infrastructure/apply?auto_deploy=false\" \\\n  -H \"Authorization: Bearer YOUR_API_TOKEN\" \\\n  -H \"Content-Type: application/yaml\" \\\n  --data-binary @ploi.yaml\n```\n\n## Response Format\n\n```json\n{\n  \"application_id\": 123,\n  \"application_name\": \"my-app\",\n  \"changes\": [\n    \"Application 'my-app' created\",\n    \"Domain 'app.example.com' added\",\n    \"Secret 'APP_KEY' created\",\n    \"Service 'database' (mysql) created\"\n  ],\n  \"structured_changes\": {\n    \"application\": { \"action\": \"created\", \"name\": \"my-app\" },\n    \"domains\": [{ \"action\": \"created\", \"domain\": \"app.example.com\" }],\n    \"secrets\": [{ \"action\": \"created\", \"key\": \"APP_KEY\" }],\n    \"services\": [{ \"action\": \"created\", \"name\": \"database\", \"type\": \"mysql\" }]\n  },\n  \"needs_deployment\": true,\n  \"auto_deploy_enabled\": true,\n  \"deployment_id\": 456,        // Only present if auto_deploy=true and deployment started\n  \"dry_run\": false,             // Indicates if this was a dry run\n  \"team\": \"My Team\"\n}\n```\n\n## Dry Run Mode\n\nWhen `dry_run=true`:\n- No database changes are made\n- No resources are created or modified\n- Response shows what would be changed with \"[DRY RUN] Would...\" prefixes\n- Deployments are never triggered, even if `auto_deploy=true`\n- Use this to validate your YAML and preview changes before applying",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gInfrastructure.AddCommand(cInfrastructureApply)

	gTeams := &cobra.Command{Use: "teams", Short: "Manage teams"}
	root.AddCommand(gTeams)
	cTeamsIndex := buildCmd(c, opSpec{
		ID:          "teams.index",
		Method:      "GET",
		PathTmpl:    "/teams",
		Use:         "list",
		Short:       "List user teams",
		Long:        "Retrieve all teams the authenticated user belongs to, including teams they own and teams\nwhere they are a member. Each team represents an isolated workspace with its own applications\nand members. The response includes aggregate counts for applications and members.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTeams.AddCommand(cTeamsIndex)

	gTeamsNetworks := &cobra.Command{Use: "networks", Short: "Manage networks"}
	gTeams.AddCommand(gTeamsNetworks)
	cTeamsNetworksDestroy := buildCmd(c, opSpec{
		ID:       "teams.networks.destroy",
		Method:   "DELETE",
		PathTmpl: "/teams/{team}/networks/{network}",
		Use:      "delete",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
			{Name: "network", Type: "integer", Required: true, Desc: "The network ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTeamsNetworks.AddCommand(cTeamsNetworksDestroy)

	cTeamsNetworksIndex := buildCmd(c, opSpec{
		ID:       "teams.networks.index",
		Method:   "GET",
		PathTmpl: "/teams/{team}/networks",
		Use:      "list",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTeamsNetworks.AddCommand(cTeamsNetworksIndex)

	gTeamsNetworksMembers := &cobra.Command{Use: "members", Short: "Manage members"}
	gTeamsNetworks.AddCommand(gTeamsNetworksMembers)
	cTeamsNetworksMembersAdd := buildCmd(c, opSpec{
		ID:       "teams.networks.members.add",
		Method:   "POST",
		PathTmpl: "/teams/{team}/networks/{network}/members",
		Use:      "add",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
			{Name: "network", Type: "integer", Required: true, Desc: "The network ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "application_id", Type: "integer", Required: true},
		},
	})
	gTeamsNetworksMembers.AddCommand(cTeamsNetworksMembersAdd)

	cTeamsNetworksMembersRemove := buildCmd(c, opSpec{
		ID:       "teams.networks.members.remove",
		Method:   "DELETE",
		PathTmpl: "/teams/{team}/networks/{network}/members",
		Use:      "remove",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
			{Name: "network", Type: "integer", Required: true, Desc: "The network ID"},
		},
		QueryParams: []paramDef{
			{Name: "application_id", Type: "integer", Required: true, Desc: ""},
		},
		BodyParams: []paramDef{},
	})
	gTeamsNetworksMembers.AddCommand(cTeamsNetworksMembersRemove)

	cTeamsNetworksShow := buildCmd(c, opSpec{
		ID:       "teams.networks.show",
		Method:   "GET",
		PathTmpl: "/teams/{team}/networks/{network}",
		Use:      "get",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
			{Name: "network", Type: "integer", Required: true, Desc: "The network ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTeamsNetworks.AddCommand(cTeamsNetworksShow)

	cTeamsNetworksStore := buildCmd(c, opSpec{
		ID:       "teams.networks.store",
		Method:   "POST",
		PathTmpl: "/teams/{team}/networks",
		Use:      "create",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "description", Type: "string", Required: false},
			{Name: "name", Type: "string", Required: true},
		},
	})
	gTeamsNetworks.AddCommand(cTeamsNetworksStore)

	cTeamsNetworksUpdate := buildCmd(c, opSpec{
		ID:       "teams.networks.update",
		Method:   "PUT",
		PathTmpl: "/teams/{team}/networks/{network}",
		Use:      "update",
		Short:    "",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
			{Name: "network", Type: "integer", Required: true, Desc: "The network ID"},
		},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "description", Type: "string", Required: false},
			{Name: "name", Type: "string", Required: false},
		},
	})
	gTeamsNetworks.AddCommand(cTeamsNetworksUpdate)

	cTeamsShow := buildCmd(c, opSpec{
		ID:       "teams.show",
		Method:   "GET",
		PathTmpl: "/teams/{team}",
		Use:      "get",
		Short:    "Get team details",
		Long:     "Retrieve detailed information about a specific team. Only accessible to team members.\nShows team information, member count, application count, and the current user's role\nwithin the team.",
		PathParams: []paramDef{
			{Name: "team", Type: "integer", Required: true, Desc: "The team ID"},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTeams.AddCommand(cTeamsShow)

	gTools := &cobra.Command{Use: "tools", Short: "Manage tools"}
	root.AddCommand(gTools)
	cToolsBalance := buildCmd(c, opSpec{
		ID:          "tools.balance",
		Method:      "GET",
		PathTmpl:    "/tools/balance",
		Use:         "balance",
		Short:       "Get PDF API credits balance",
		Long:        "Retrieve the current team's credit balance and PDF API usage statistics.\nThis endpoint does not consume credits.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTools.AddCommand(cToolsBalance)

	cToolsHealth := buildCmd(c, opSpec{
		ID:          "tools.health",
		Method:      "GET",
		PathTmpl:    "/tools/health",
		Use:         "health",
		Short:       "Check PDF service health",
		Long:        "Verify that the PDF conversion service is available and operational.\nThis endpoint does not consume credits.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gTools.AddCommand(cToolsHealth)

	gToolsPdf := &cobra.Command{Use: "pdf", Short: "Manage pdf"}
	gTools.AddCommand(gToolsPdf)
	gToolsPdfConvert := &cobra.Command{Use: "convert", Short: "Manage convert"}
	gToolsPdf.AddCommand(gToolsPdfConvert)
	cToolsPdfConvertHtml := buildCmd(c, opSpec{
		ID:          "tools.pdf.convert.html",
		Method:      "POST",
		PathTmpl:    "/tools/pdf/convert/html",
		Use:         "html",
		Short:       "Convert HTML to PDF",
		Long:        "Convert HTML content to PDF format. Supports inline CSS and images (base64 encoded).\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "html", Type: "string", Required: true},
			{Name: "landscape", Type: "boolean", Required: false},
			{Name: "marginBottom", Type: "string", Required: false},
			{Name: "marginLeft", Type: "string", Required: false},
			{Name: "marginRight", Type: "string", Required: false},
			{Name: "marginTop", Type: "string", Required: false},
			{Name: "paperHeight", Type: "string", Required: false},
			{Name: "paperWidth", Type: "string", Required: false},
			{Name: "printBackground", Type: "boolean", Required: false},
			{Name: "scale", Type: "number", Required: false},
		},
	})
	gToolsPdfConvert.AddCommand(cToolsPdfConvertHtml)

	cToolsPdfConvertMarkdown := buildCmd(c, opSpec{
		ID:          "tools.pdf.convert.markdown",
		Method:      "POST",
		PathTmpl:    "/tools/pdf/convert/markdown",
		Use:         "markdown",
		Short:       "Convert Markdown to PDF",
		Long:        "Convert Markdown content to PDF format with automatic styling.\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "html", Type: "string", Required: false},
			{Name: "markdown", Type: "string", Required: true},
		},
	})
	gToolsPdfConvert.AddCommand(cToolsPdfConvertMarkdown)

	cToolsPdfConvertOffice := buildCmd(c, opSpec{
		ID:          "tools.pdf.convert.office",
		Method:      "POST",
		PathTmpl:    "/tools/pdf/convert/office",
		Use:         "office",
		Short:       "Convert Office document to PDF",
		Long:        "Convert Microsoft Office documents (Word, Excel, PowerPoint) and other formats to PDF.\nSupported formats: .docx, .doc, .xlsx, .xls, .pptx, .ppt, .odt, .ods, .odp, .rtf, .txt, .csv\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gToolsPdfConvert.AddCommand(cToolsPdfConvertOffice)

	cToolsPdfConvertUrl := buildCmd(c, opSpec{
		ID:          "tools.pdf.convert.url",
		Method:      "POST",
		PathTmpl:    "/tools/pdf/convert/url",
		Use:         "url",
		Short:       "Convert URL to PDF",
		Long:        "Convert a web page to PDF format. The URL must be publicly accessible.\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "landscape", Type: "boolean", Required: false},
			{Name: "marginBottom", Type: "string", Required: false},
			{Name: "marginLeft", Type: "string", Required: false},
			{Name: "marginRight", Type: "string", Required: false},
			{Name: "marginTop", Type: "string", Required: false},
			{Name: "paperHeight", Type: "string", Required: false},
			{Name: "paperWidth", Type: "string", Required: false},
			{Name: "printBackground", Type: "boolean", Required: false},
			{Name: "scale", Type: "number", Required: false},
			{Name: "url", Type: "string", Required: true},
		},
	})
	gToolsPdfConvert.AddCommand(cToolsPdfConvertUrl)

	cToolsPdfMerge := buildCmd(c, opSpec{
		ID:          "tools.pdf.merge",
		Method:      "POST",
		PathTmpl:    "/tools/pdf/merge",
		Use:         "merge",
		Short:       "Merge multiple PDFs",
		Long:        "Combine multiple PDF files into a single PDF document.\nFiles are merged in the order they are provided.\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gToolsPdf.AddCommand(cToolsPdfMerge)

	gToolsScreenshot := &cobra.Command{Use: "screenshot", Short: "Manage screenshot"}
	gTools.AddCommand(gToolsScreenshot)
	cToolsScreenshotHtml := buildCmd(c, opSpec{
		ID:          "tools.screenshot.html",
		Method:      "POST",
		PathTmpl:    "/tools/screenshot/html",
		Use:         "html",
		Short:       "Take HTML screenshot",
		Long:        "Capture a screenshot of HTML content. Returns PNG or JPEG image.\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "format", Type: "string", Required: false},
			{Name: "height", Type: "integer", Required: false},
			{Name: "html", Type: "string", Required: true},
			{Name: "width", Type: "integer", Required: false},
		},
	})
	gToolsScreenshot.AddCommand(cToolsScreenshotHtml)

	cToolsScreenshotUrl := buildCmd(c, opSpec{
		ID:          "tools.screenshot.url",
		Method:      "POST",
		PathTmpl:    "/tools/screenshot/url",
		Use:         "url",
		Short:       "Take URL screenshot",
		Long:        "Capture a screenshot of a web page. Returns PNG or JPEG image.\nConsumes 1 credit per request.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "format", Type: "string", Required: false},
			{Name: "height", Type: "integer", Required: false},
			{Name: "url", Type: "string", Required: true},
			{Name: "width", Type: "integer", Required: false},
		},
	})
	gToolsScreenshot.AddCommand(cToolsScreenshotUrl)

	gUser := &cobra.Command{Use: "user", Short: "Manage user"}
	root.AddCommand(gUser)
	cUserShow := buildCmd(c, opSpec{
		ID:          "user.show",
		Method:      "GET",
		PathTmpl:    "/user",
		Use:         "get",
		Short:       "Get current user",
		Long:        "Retrieve detailed information about the authenticated user including their profile,\ncurrent team context, and all teams they belong to. The avatar URL is generated\nusing Gravatar based on the user's email address.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gUser.AddCommand(cUserShow)

	cUserSwitchTeam := buildCmd(c, opSpec{
		ID:          "user.switch-team",
		Method:      "PUT",
		PathTmpl:    "/user/team",
		Use:         "switch-team",
		Short:       "Switch current team",
		Long:        "Switch the authenticated user's current team context. All subsequent API requests\nwill operate within the context of the selected team. The user must be a member\nof the team to switch to it.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "team_id", Type: "integer", Required: true},
		},
	})
	gUser.AddCommand(cUserSwitchTeam)

	gUserTokens := &cobra.Command{Use: "tokens", Short: "Manage tokens"}
	gUser.AddCommand(gUserTokens)
	cUserTokensDestroy := buildCmd(c, opSpec{
		ID:       "user.tokens.destroy",
		Method:   "DELETE",
		PathTmpl: "/user/tokens/{tokenId}",
		Use:      "delete",
		Short:    "Revoke API token",
		Long:     "Permanently revoke a personal access token. This immediately invalidates the token\nand any requests using it will fail with 401 Unauthorized. This action cannot be\nundone - you'll need to create a new token if needed.",
		PathParams: []paramDef{
			{Name: "tokenId", Type: "string", Required: true, Desc: ""},
		},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gUserTokens.AddCommand(cUserTokensDestroy)

	cUserTokensIndex := buildCmd(c, opSpec{
		ID:          "user.tokens.index",
		Method:      "GET",
		PathTmpl:    "/user/tokens",
		Use:         "list",
		Short:       "List API tokens",
		Long:        "Retrieve all active personal access tokens for the authenticated user. Tokens are\nused for API authentication and can have specific scopes or full access. The actual\ntoken values are not included in the response for security reasons.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams:  []paramDef{},
	})
	gUserTokens.AddCommand(cUserTokensIndex)

	cUserTokensStore := buildCmd(c, opSpec{
		ID:          "user.tokens.store",
		Method:      "POST",
		PathTmpl:    "/user/tokens",
		Use:         "create",
		Short:       "Create API token",
		Long:        "Create a new personal access token for API authentication. The token value is only\nshown once during creation and cannot be retrieved again. Store it securely. Tokens\ncan have specific scopes to limit permissions or use wildcard (*) for full access.\nPassword verification is required for security.",
		PathParams:  []paramDef{},
		QueryParams: []paramDef{},
		BodyParams: []paramDef{
			{Name: "expires_at", Type: "string", Required: false},
			{Name: "name", Type: "string", Required: true},
			{Name: "password", Type: "string", Required: true},
			{Name: "scopes", Type: "array:string", Required: false},
		},
	})
	gUserTokens.AddCommand(cUserTokensStore)

}
