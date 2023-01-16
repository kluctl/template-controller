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
<p>Interval is the interval at which to query the Gitlab API.
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
<p>Interval is the interval at which to query the Gitlab API.
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
<p>Branch to filter for. Can also be a regex.</p>
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
<p>Authentication token reference.</p>
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
<a href="#templates.kluctl.io/v1alpha1.GithubCommentSpec">GithubCommentSpec</a>, 
<a href="#templates.kluctl.io/v1alpha1.PullRequestRefHolder">PullRequestRefHolder</a>)
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
int
</em>
</td>
<td>
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
<a href="#templates.kluctl.io/v1alpha1.GitlabCommentSpec">GitlabCommentSpec</a>, 
<a href="#templates.kluctl.io/v1alpha1.PullRequestRefHolder">PullRequestRefHolder</a>)
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
int
</em>
</td>
<td>
<p>The merge request id</p>
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
string
</em>
</td>
<td>
<p>GitLab project to scan. Required.</p>
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
<p>The GitLab API URL to talk to. If blank, uses <a href="https://gitlab.com/">https://gitlab.com/</a>.</p>
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
<p>Authentication token reference.</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.Handler">Handler
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectHandlerSpec">ObjectHandlerSpec</a>)
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
<code>pullRequestComment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommentReporter">
PullRequestCommentReporter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>pullRequestApprove</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestApproveReporter">
PullRequestApproveReporter
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>pullRequestCommand</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandler">
PullRequestCommandHandler
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
<h3 id="templates.kluctl.io/v1alpha1.HandlerStatus">HandlerStatus
</h3>
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
<code>key</code><br>
<em>
string
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
<tr>
<td>
<code>pullRequestComment</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommentReporterStatus">
PullRequestCommentReporterStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>pullRequestApprove</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestApproveReporterStatus">
PullRequestApproveReporterStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>pullRequestCommand</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerStatus">
PullRequestCommandHandlerStatus
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
<a href="#templates.kluctl.io/v1alpha1.MatrixEntryObject">
MatrixEntryObject
</a>
</em>
</td>
<td>
<em>(Optional)</em>
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
<tr>
<td>
<code>expandLists</code><br>
<em>
bool
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
<h3 id="templates.kluctl.io/v1alpha1.ObjectHandler">ObjectHandler
</h3>
<p>ObjectHandler is the Schema for the objecthandlers API</p>
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
<a href="#templates.kluctl.io/v1alpha1.ObjectHandlerSpec">
ObjectHandlerSpec
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
<code>forObject</code><br>
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
<code>handlers</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.Handler">
[]Handler
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.ObjectHandlerStatus">
ObjectHandlerStatus
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
<h3 id="templates.kluctl.io/v1alpha1.ObjectHandlerSpec">ObjectHandlerSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectHandler">ObjectHandler</a>)
</p>
<p>ObjectHandlerSpec defines the desired state of ObjectHandler</p>
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
<code>forObject</code><br>
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
<code>handlers</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.Handler">
[]Handler
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
<h3 id="templates.kluctl.io/v1alpha1.ObjectHandlerStatus">ObjectHandlerStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.ObjectHandler">ObjectHandler</a>)
</p>
<p>ObjectHandlerStatus defines the observed state of ObjectHandler</p>
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
<code>handlerStatus</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.*github.com/kluctl/template-controller/api/v1alpha1.HandlerStatus">
[]*github.com/kluctl/template-controller/api/v1alpha1.HandlerStatus
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
<h3 id="templates.kluctl.io/v1alpha1.ObjectRef">ObjectRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.AppliedResourceInfo">AppliedResourceInfo</a>, 
<a href="#templates.kluctl.io/v1alpha1.MatrixEntryObject">MatrixEntryObject</a>, 
<a href="#templates.kluctl.io/v1alpha1.ObjectHandlerSpec">ObjectHandlerSpec</a>, 
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
<p>The name of the Kubernetes service account to impersonate
when reconciling this ObjectTemplate. If omitted, the &ldquo;default&rdquo; service account is used.</p>
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
</td>
</tr>
<tr>
<td>
<code>matrix</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.*github.com/kluctl/template-controller/api/v1alpha1.MatrixEntry">
[]*github.com/kluctl/template-controller/api/v1alpha1.MatrixEntry
</a>
</em>
</td>
<td>
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
<p>The name of the Kubernetes service account to impersonate
when reconciling this ObjectTemplate. If omitted, the &ldquo;default&rdquo; service account is used.</p>
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
</td>
</tr>
<tr>
<td>
<code>matrix</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.*github.com/kluctl/template-controller/api/v1alpha1.MatrixEntry">
[]*github.com/kluctl/template-controller/api/v1alpha1.MatrixEntry
</a>
</em>
</td>
<td>
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestApproveReporter">PullRequestApproveReporter
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.Handler">Handler</a>)
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
<code>PullRequestRefHolder</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestRefHolder">
PullRequestRefHolder
</a>
</em>
</td>
<td>
<p>
(Members of <code>PullRequestRefHolder</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>missingReadyConditionIsError</code><br>
<em>
bool
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestApproveReporterStatus">PullRequestApproveReporterStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.HandlerStatus">HandlerStatus</a>)
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
<code>approved</code><br>
<em>
bool
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommandHandler">PullRequestCommandHandler
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.Handler">Handler</a>)
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
<code>PullRequestRefHolder</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestRefHolder">
PullRequestRefHolder
</a>
</em>
</td>
<td>
<p>
(Members of <code>PullRequestRefHolder</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>postHelpComment</code><br>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>commands</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerCommandSpec">
[]PullRequestCommandHandlerCommandSpec
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommandHandlerActionAnnotateSpec">PullRequestCommandHandlerActionAnnotateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerActionSpec">PullRequestCommandHandlerActionSpec</a>)
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
<code>annotation</code><br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>value</code><br>
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommandHandlerActionSpec">PullRequestCommandHandlerActionSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerCommandSpec">PullRequestCommandHandlerCommandSpec</a>)
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
<code>annotate</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerActionAnnotateSpec">
PullRequestCommandHandlerActionAnnotateSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>jsonPatch</code><br>
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommandHandlerCommandSpec">PullRequestCommandHandlerCommandSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandler">PullRequestCommandHandler</a>)
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
<code>description</code><br>
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
<code>actions</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandlerActionSpec">
[]PullRequestCommandHandlerActionSpec
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommandHandlerStatus">PullRequestCommandHandlerStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.HandlerStatus">HandlerStatus</a>)
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
<code>lastProcessedCommentTime</code><br>
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
<code>helpNoteId</code><br>
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
<code>helpNoteBodyHash</code><br>
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommentReporter">PullRequestCommentReporter
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.Handler">Handler</a>)
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
<code>PullRequestRefHolder</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestRefHolder">
PullRequestRefHolder
</a>
</em>
</td>
<td>
<p>
(Members of <code>PullRequestRefHolder</code> are embedded into this type.)
</p>
</td>
</tr>
</tbody>
</table>
</div>
</div>
<h3 id="templates.kluctl.io/v1alpha1.PullRequestCommentReporterStatus">PullRequestCommentReporterStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.HandlerStatus">HandlerStatus</a>)
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
<code>lastPostedStatusHash</code><br>
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
<code>noteId</code><br>
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
<h3 id="templates.kluctl.io/v1alpha1.PullRequestRefHolder">PullRequestRefHolder
</h3>
<p>
(<em>Appears on:</em>
<a href="#templates.kluctl.io/v1alpha1.PullRequestApproveReporter">PullRequestApproveReporter</a>, 
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommandHandler">PullRequestCommandHandler</a>, 
<a href="#templates.kluctl.io/v1alpha1.PullRequestCommentReporter">PullRequestCommentReporter</a>)
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
<code>gitlab</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.GitlabMergeRequestRef">
GitlabMergeRequestRef
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
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
<p>The name of the Kubernetes service account to impersonate
when reconciling this TextTemplate. If omitted, the &ldquo;default&rdquo; service account is used.</p>
</td>
</tr>
<tr>
<td>
<code>inputs</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.*github.com/kluctl/template-controller/api/v1alpha1.TextTemplateInput">
[]*github.com/kluctl/template-controller/api/v1alpha1.TextTemplateInput
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
<p>The name of the Kubernetes service account to impersonate
when reconciling this TextTemplate. If omitted, the &ldquo;default&rdquo; service account is used.</p>
</td>
</tr>
<tr>
<td>
<code>inputs</code><br>
<em>
<a href="#templates.kluctl.io/v1alpha1.*github.com/kluctl/template-controller/api/v1alpha1.TextTemplateInput">
[]*github.com/kluctl/template-controller/api/v1alpha1.TextTemplateInput
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
