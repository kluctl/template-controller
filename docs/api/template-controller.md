<h1>Template Controller API reference</h1>
<p>Packages:</p>
<ul class="simple">
<li>
<a href="#templates.kluctl.io%2fv1alpha1">templates.kluctl.io/v1alpha1</a>
</li>
</ul>
<h2 id="templates.kluctl.io/v1alpha1">templates.kluctl.io/v1alpha1</h2>
<p>Package v1alpha1 contains API Schema definitions for the templates.kluctl.io v1alpha1 API group.</p>
Resource Types:
<ul class="simple"></ul>
<h3 id="templates.kluctl.io/v1alpha1.AppliedResourceInfo">AppliedResourceInfo
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplateStatus">ObjectTemplateStatus</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectRef">
ObjectRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>success</code><br>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>error</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.CommentSourceSpec">CommentSourceSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.CommentSpec">CommentSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>text</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Text specifies a raw text comment.</p>
</td>
</tr>
<tr>
<td>
<code>configMap</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ConfigMapRef">
ConfigMapRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ConfigMap specifies a ConfigMap and a key to load the source content from</p>
</td>
</tr>
<tr>
<td>
<code>textTemplate</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TextTemplate specifies a TextTemplate to load the source content from</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.CommentSpec">CommentSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubCommentSpec">GithubCommentSpec</a>, 
<a href="#templates.kluctl.io/v1alpha1.GitlabCommentSpec">GitlabCommentSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>id</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Id specifies the identifier to be used by the controller when it needs to find the actual comment when it does
not know the internal id. This Id is written into the comment inside a comment, so that a simple text search
can reveal the comment</p>
</td>
</tr>
<tr>
<td>
<code>source</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.CommentSourceSpec">
CommentSourceSpec
</a>
</em>
</td>
<td>
<p>Source specifies the source content for the comment. Different source types are supported:
Text, ConfigMap and TextTemplate</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ConfigMapRef">ConfigMapRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.CommentSourceSpec">CommentSourceSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>key</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitFile">GitFile
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorSpec">GitProjectorSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>glob</code><br>
<em>
string
</em>
</td>
<td>
<p>Glob specifies a glob to use for filename matching.</p>
</td>
</tr>
<tr>
<td>
<code>parseYaml</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>ParseYaml enables YAML parsing of matching files. The result is then available as <code>parsed</code> in the result for
the corresponding result file</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitProjector">GitProjector
</h3>
<p>GitProjector is the Schema for the gitprojectors API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorSpec">
GitProjectorSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to scan the Git repository
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
<tr>
<td>
<code>url</code><br>
<em>
string
</em>
</td>
<td>
<p>URL specifies the Git url to scan and project</p>
</td>
</tr>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitRef">
GitRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reference specifies the Git branch, tag or commit to scan. Branches and tags can contain regular expressions</p>
</td>
</tr>
<tr>
<td>
<code>files</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitFile">
[]GitFile
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Files specifies the list of files to include in the projection</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretRefs specifies a Secret use for Git authentication. The contents of the secret must conform to:
<a href="https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/#git-authentication">https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/#git-authentication</a></p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorStatus">
GitProjectorStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitProjectorResult">GitProjectorResult
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorStatus">GitProjectorStatus</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitRef">
GitRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>files</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorResultFile">
[]GitProjectorResultFile
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitProjectorResultFile">GitProjectorResultFile
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorResult">GitProjectorResult</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>raw</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>parsed</code><br>
<em>
[]*k8s.io/apimachinery/pkg/runtime.RawExtension
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitProjectorSpec">GitProjectorSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjector">GitProjector</a>)
</p>
<p>GitProjectorSpec defines the desired state of GitProjector</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to scan the Git repository
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
<tr>
<td>
<code>url</code><br>
<em>
string
</em>
</td>
<td>
<p>URL specifies the Git url to scan and project</p>
</td>
</tr>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitRef">
GitRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Reference specifies the Git branch, tag or commit to scan. Branches and tags can contain regular expressions</p>
</td>
</tr>
<tr>
<td>
<code>files</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitFile">
[]GitFile
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Files specifies the list of files to include in the projection</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.LocalObjectReference">
LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretRefs specifies a Secret use for Git authentication. The contents of the secret must conform to:
<a href="https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/#git-authentication">https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/#git-authentication</a></p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitProjectorStatus">GitProjectorStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjector">GitProjector</a>)
</p>
<p>GitProjectorStatus defines the observed state of GitProjector</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>allRefsHash</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>result</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorResult">
[]GitProjectorResult
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitRef">GitRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitProjectorResult">GitProjectorResult</a>, 
<a href="#templates.kluctl.io/v1alpha1.GitProjectorSpec">GitProjectorSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>branch</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Branch to filter for. Can also be a regex.</p>
</td>
</tr>
<tr>
<td>
<code>tag</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tag to filter for. Can also be a regex.</p>
</td>
</tr>
<tr>
<td>
<code>commit</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Commit SHA to check out, takes precedence over all reference fields.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GithubComment">GithubComment
</h3>
<p>GithubComment is the Schema for the githubcomments API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubCommentSpec">
GithubCommentSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>github</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubPullRequestRef">
GithubPullRequestRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>comment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.CommentSpec">
CommentSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubCommentStatus">
GithubCommentStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GithubCommentSpec">GithubCommentSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubComment">GithubComment</a>)
</p>
<p>GithubCommentSpec defines the desired state of GithubComment</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>github</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubPullRequestRef">
GithubPullRequestRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>comment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.CommentSpec">
CommentSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GithubCommentStatus">GithubCommentStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubComment">GithubComment</a>)
</p>
<p>GithubCommentStatus defines the observed state of GithubComment</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>commentId</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>lastPostedBodyHash</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GithubProject">GithubProject
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubPullRequestRef">GithubPullRequestRef</a>, 
<a href="#templates.kluctl.io/v1alpha1.ListGithubPullRequestsSpec">ListGithubPullRequestsSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>owner</code><br>
<em>
string
</em>
</td>
<td>
<p>Owner specifies the GitHub user or organisation that owns the repository</p>
</td>
</tr>
<tr>
<td>
<code>repo</code><br>
<em>
string
</em>
</td>
<td>
<p>Repo specifies the repository name.</p>
</td>
</tr>
<tr>
<td>
<code>tokenRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TokenRef specifies a secret and key to load the GitHub API token from</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GithubPullRequestRef">GithubPullRequestRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubCommentSpec">GithubCommentSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>GithubProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubProject">
GithubProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GithubProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>pullRequestId</code><br>
<em>
k8s.io/apimachinery/pkg/util/intstr.IntOrString
</em>
</td>
<td>
<p>PullRequestId specifies the pull request ID.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitlabComment">GitlabComment
</h3>
<p>GitlabComment is the Schema for the gitlabcomments API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabCommentSpec">
GitlabCommentSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>gitlab</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabMergeRequestRef">
GitlabMergeRequestRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>comment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.CommentSpec">
CommentSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabCommentStatus">
GitlabCommentStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitlabCommentSpec">GitlabCommentSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitlabComment">GitlabComment</a>)
</p>
<p>GitlabCommentSpec defines the desired state of GitlabComment</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>gitlab</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabMergeRequestRef">
GitlabMergeRequestRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>comment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.CommentSpec">
CommentSpec
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitlabCommentStatus">GitlabCommentStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitlabComment">GitlabComment</a>)
</p>
<p>GitlabCommentStatus defines the observed state of GitlabComment</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>noteId</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>lastPostedBodyHash</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitlabMergeRequestRef">GitlabMergeRequestRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitlabCommentSpec">GitlabCommentSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>GitlabProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabProject">
GitlabProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GitlabProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>mergeRequestId</code><br>
<em>
k8s.io/apimachinery/pkg/util/intstr.IntOrString
</em>
</td>
<td>
<p>MergeRequestId specifies the Gitlab merge request internal ID</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.GitlabProject">GitlabProject
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GitlabMergeRequestRef">GitlabMergeRequestRef</a>, 
<a href="#templates.kluctl.io/v1alpha1.ListGitlabMergeRequestsSpec">ListGitlabMergeRequestsSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>project</code><br>
<em>
k8s.io/apimachinery/pkg/util/intstr.IntOrString
</em>
</td>
<td>
<p>Project specifies the Gitlab group and project (separated by slash) to
use, or the numeric project id</p>
</td>
</tr>
<tr>
<td>
<code>api</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>API specifies the GitLab API URL to talk to.
If blank, uses <a href="https://gitlab.com/">https://gitlab.com/</a>.</p>
</td>
</tr>
<tr>
<td>
<code>tokenRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.SecretRef">
SecretRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>TokenRef specifies a secret and key to load the Gitlab API token from</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGithubPullRequests">ListGithubPullRequests
</h3>
<p>ListGithubPullRequests is the Schema for the listgithubpullrequests API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ListGithubPullRequestsSpec">
ListGithubPullRequestsSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to query the Gitlab API.
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>GithubProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubProject">
GithubProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GithubProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>head</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Head specifies the head to filter for</p>
</td>
</tr>
<tr>
<td>
<code>base</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Base specifies the base to filter for</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels is used to filter the PRs that you want to target</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>State is an additional PR filter to get only those with a certain state. Default: &ldquo;all&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>limit</code><br>
<em>
int
</em>
</td>
<td>
<p>Limit limits the maximum number of pull requests to fetch. Defaults to 100</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ListGithubPullRequestsStatus">
ListGithubPullRequestsStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGithubPullRequestsSpec">ListGithubPullRequestsSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ListGithubPullRequests">ListGithubPullRequests</a>)
</p>
<p>ListGithubPullRequestsSpec defines the desired state of ListGithubPullRequests</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to query the Gitlab API.
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>GithubProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GithubProject">
GithubProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GithubProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>head</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Head specifies the head to filter for</p>
</td>
</tr>
<tr>
<td>
<code>base</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Base specifies the base to filter for</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels is used to filter the PRs that you want to target</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>State is an additional PR filter to get only those with a certain state. Default: &ldquo;all&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>limit</code><br>
<em>
int
</em>
</td>
<td>
<p>Limit limits the maximum number of pull requests to fetch. Defaults to 100</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGithubPullRequestsStatus">ListGithubPullRequestsStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ListGithubPullRequests">ListGithubPullRequests</a>)
</p>
<p>ListGithubPullRequestsStatus defines the observed state of ListGithubPullRequests</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>pullRequests</code><br>
<em>
[]k8s.io/apimachinery/pkg/runtime.RawExtension
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGitlabMergeRequests">ListGitlabMergeRequests
</h3>
<p>ListGitlabMergeRequests is the Schema for the listgitlabmergerequests API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ListGitlabMergeRequestsSpec">
ListGitlabMergeRequestsSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to query the Gitlab API.
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>GitlabProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabProject">
GitlabProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GitlabProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>targetBranch</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TargetBranch specifies the target branch to filter for</p>
</td>
</tr>
<tr>
<td>
<code>sourceBranch</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>labels</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels is used to filter the MRs that you want to target</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>State is an additional MRs filter to get only those with a certain state. Default: &ldquo;all&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>limit</code><br>
<em>
int
</em>
</td>
<td>
<p>Limit limits the maximum number of merge requests to fetch. Defaults to 100</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ListGitlabMergeRequestsStatus">
ListGitlabMergeRequestsStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGitlabMergeRequestsSpec">ListGitlabMergeRequestsSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ListGitlabMergeRequests">ListGitlabMergeRequests</a>)
</p>
<p>ListGitlabMergeRequestsSpec defines the desired state of ListGitlabMergeRequests</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Interval is the interval at which to query the Gitlab API.
Defaults to 5m.</p>
</td>
</tr>
<tr>
<td>
<code>GitlabProject</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabProject">
GitlabProject
</a>
</em>
</td>
<td>
<p>
(Members of <code>GitlabProject</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>targetBranch</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>TargetBranch specifies the target branch to filter for</p>
</td>
</tr>
<tr>
<td>
<code>sourceBranch</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>labels</code><br>
<em>
[]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Labels is used to filter the MRs that you want to target</p>
</td>
</tr>
<tr>
<td>
<code>state</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>State is an additional MRs filter to get only those with a certain state. Default: &ldquo;all&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>limit</code><br>
<em>
int
</em>
</td>
<td>
<p>Limit limits the maximum number of merge requests to fetch. Defaults to 100</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ListGitlabMergeRequestsStatus">ListGitlabMergeRequestsStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ListGitlabMergeRequests">ListGitlabMergeRequests</a>)
</p>
<p>ListGitlabMergeRequestsStatus defines the observed state of ListGitlabMergeRequests</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>mergeRequests</code><br>
<em>
[]k8s.io/apimachinery/pkg/runtime.RawExtension
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.LocalObjectReference">LocalObjectReference
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.CommentSourceSpec">CommentSourceSpec</a>, 
<a href="#templates.kluctl.io/v1alpha1.GitProjectorSpec">GitProjectorSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
<p>Name of the referent.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.MatrixEntry">MatrixEntry
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplateSpec">ObjectTemplateSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
<p>Name specifies the name this matrix input is available while rendering templates</p>
</td>
</tr>
<tr>
<td>
<code>object</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.MatrixEntryObject">
MatrixEntryObject
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Object specifies an object to load and make available while rendering templates. The object can be accessed
through the name specified above. The service account used by the ObjectTemplate must have proper permissions
to get this object</p>
</td>
</tr>
<tr>
<td>
<code>list</code><br>
<em>
[]k8s.io/apimachinery/pkg/runtime.RawExtension
</em>
</td>
<td>
<em>(Optional)</em>
<p>List specifies a list of plain YAML values which are made available while rendering templates. The list can be
accessed through the name specified above</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.MatrixEntryObject">MatrixEntryObject
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.MatrixEntry">MatrixEntry</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectRef">
ObjectRef
</a>
</em>
</td>
<td>
<p>Ref specifies the apiVersion, kind, namespace and name of the object to load. The service account used by the
ObjectTemplate must have proper permissions to get this object</p>
</td>
</tr>
<tr>
<td>
<code>jsonPath</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>JsonPath optionally specifies a sub-field to load. When specified, the sub-field (and not the whole object)
is made available while rendering templates</p>
</td>
</tr>
<tr>
<td>
<code>expandLists</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExpandLists enables optional expanding of list. Expanding means, that each list entry is interpreted as
individual matrix input instead of interpreting the whole list as one matrix input. This feature is only useful
when used in combination with <code>jsonPath</code></p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ObjectRef">ObjectRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.AppliedResourceInfo">AppliedResourceInfo</a>, 
<a href="#templates.kluctl.io/v1alpha1.MatrixEntryObject">MatrixEntryObject</a>, 
<a href="#templates.kluctl.io/v1alpha1.TextTemplateInputObject">TextTemplateInputObject</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kind</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ObjectTemplate">ObjectTemplate
</h3>
<p>ObjectTemplate is the Schema for the objecttemplates API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplateSpec">
ObjectTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName specifies the name of the Kubernetes service account to impersonate
when reconciling this ObjectTemplate. If omitted, the &ldquo;default&rdquo; service account is used</p>
</td>
</tr>
<tr>
<td>
<code>prune</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Prune enables pruning of previously created objects when these disappear from the list of rendered objects</p>
</td>
</tr>
<tr>
<td>
<code>matrix</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.MatrixEntry">
[]MatrixEntry
</a>
</em>
</td>
<td>
<p>Matrix specifies the input matrix</p>
</td>
</tr>
<tr>
<td>
<code>templates</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.Template">
[]Template
</a>
</em>
</td>
<td>
<p>Templates specifies a list of templates to render and deploy</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplateStatus">
ObjectTemplateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ObjectTemplateSpec">ObjectTemplateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplate">ObjectTemplate</a>)
</p>
<p>ObjectTemplateSpec defines the desired state of ObjectTemplate</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>interval</code><br>
<em>
<a href="https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName specifies the name of the Kubernetes service account to impersonate
when reconciling this ObjectTemplate. If omitted, the &ldquo;default&rdquo; service account is used</p>
</td>
</tr>
<tr>
<td>
<code>prune</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Prune enables pruning of previously created objects when these disappear from the list of rendered objects</p>
</td>
</tr>
<tr>
<td>
<code>matrix</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.MatrixEntry">
[]MatrixEntry
</a>
</em>
</td>
<td>
<p>Matrix specifies the input matrix</p>
</td>
</tr>
<tr>
<td>
<code>templates</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.Template">
[]Template
</a>
</em>
</td>
<td>
<p>Templates specifies a list of templates to render and deploy</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.ObjectTemplateStatus">ObjectTemplateStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplate">ObjectTemplate</a>)
</p>
<p>ObjectTemplateStatus defines the observed state of ObjectTemplate</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>appliedResources</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.AppliedResourceInfo">
[]AppliedResourceInfo
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.SecretRef">SecretRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.GithubProject">GithubProject</a>, 
<a href="#templates.kluctl.io/v1alpha1.GitlabProject">GitlabProject</a>)
</p>
<p>Utility struct for a reference to a secret key.</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>secretName</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>key</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.Template">Template
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectTemplateSpec">ObjectTemplateSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>object</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#unstructured-unstructured-v1">
Kubernetes meta/v1/unstructured.Unstructured
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Object specifies a structured object in YAML form. Each field value is rendered independently.</p>
</td>
</tr>
<tr>
<td>
<code>raw</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Raw specifies a raw string to be interpreted/parsed as YAML. The whole string is rendered in one go, allowing to
use advanced Jinja2 control structures. Raw object might also be required when a templated value must not be
interpreted as a string (which would be done in Object).</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TemplateRef">TemplateRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateSpec">TextTemplateSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>configMap</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TemplateRefConfigMap">
TemplateRefConfigMap
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TemplateRefConfigMap">TemplateRefConfigMap
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TemplateRef">TemplateRef</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>namespace</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>key</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TextTemplate">TextTemplate
</h3>
<p>TextTemplate is the Schema for the texttemplates API</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateSpec">
TextTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName specifies the name of the Kubernetes service account to impersonate
when reconciling this TextTemplate. If omitted, the &ldquo;default&rdquo; service account is used</p>
</td>
</tr>
<tr>
<td>
<code>inputs</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateInput">
[]TextTemplateInput
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>template</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>templateRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TemplateRef">
TemplateRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateStatus">
TextTemplateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TextTemplateInput">TextTemplateInput
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateSpec">TextTemplateSpec</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>object</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateInputObject">
TextTemplateInputObject
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TextTemplateInputObject">TextTemplateInputObject
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateInput">TextTemplateInput</a>)
</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ref</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectRef">
ObjectRef
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>jsonPath</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TextTemplateSpec">TextTemplateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplate">TextTemplate</a>)
</p>
<p>TextTemplateSpec defines the desired state of TextTemplate</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>suspend</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Suspend can be used to suspend the reconciliation of this object.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccountName</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceAccountName specifies the name of the Kubernetes service account to impersonate
when reconciling this TextTemplate. If omitted, the &ldquo;default&rdquo; service account is used</p>
</td>
</tr>
<tr>
<td>
<code>inputs</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplateInput">
[]TextTemplateInput
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>template</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>templateRef</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.TemplateRef">
TemplateRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.TextTemplateStatus">TextTemplateStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.TextTemplate">TextTemplate</a>)
</p>
<p>TextTemplateStatus defines the observed state of TextTemplate</p>
<div class="md-typeset__scrollwrap">
<div class="md-typeset__table">
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code><br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#condition-v1-meta">
[]Kubernetes meta/v1.Condition
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>result</code><br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<div class="admonition note">
<p class="last">This page was automatically generated with <code>gen-crd-api-reference-docs</code></p>
</div>
