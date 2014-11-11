// © 2014 Michael Stapelberg
// vim:ts=4:sw=4:et
// Opens a WebSocket connection to Debian Code Search to send and receive
// search results almost instantanously.

// NB: All of these constants needs to match those in cmd/dcs-web/querymanager.go
var packagesPerPage = 5;
var resultsPerPackage = 2;

var animationFallback;
var showConnectProgress;
var connection = new ReconnectingWebSocket('ws://' + window.location.hostname + ':28080/instantws');
var searchterm;

// fatal (bool): Whether all ongoing operations should be cancelled.
//
// permanent (bool): Whether this message will be displayed permanently (e.g.
// “search results incomplete” vs. “trying to reconnect in 3s…”)
//
// unique_id (string): If non-null, only one message of this type will be
// displayed. Can be used to display only one notification about incomplete
// search results, regardless of how many backends the server returns as
// unhealthy.
//
// message (string): The human-readable error message.
function error(fatal, permanent, unique_id, message) {
    if (unique_id !== null && $('#errors div[data-uniqueid=' + unique_id + ']').size() > 0) {
        return;
    }
    if (fatal) {
        progress(100, false, 'Error: ' + message);
    }

    var div = $('<div class="alert alert-' + (permanent ? 'danger' : 'warning') + '" role="alert"></div>');
    if (unique_id !== null) {
        div.attr('data-uniqueid', unique_id);
    }
    div.text(message);
    $('#errors').append(div);
    return div;
}

// Setting percentage to 0 means initializing the progressbar. To display some
// sort of progress to the user, we’ll set it to 10%, so any actual progress
// that is communicated from the server will need to be ≥ 10%.
//
// Setting temporary to true will reset the text to the last non-temporary text
// upon completion (which is a call with percentage == 100).
function progress(percentage, temporary, text) {
    if (percentage == 0) {
        $('#progressbar span').text(text);
        $('#progressbar .progress-bar').css('width', '10%');
        $('#progressbar .progress-bar').addClass('progress-active');
        $('#progressbar').show();
    } else {
        if (text !== null) {
            $('#progressbar span').text(text);
            if (!temporary) {
                $('#progressbar').data('old-text', text);
            }
        }
        $('#progressbar .progress-bar').css('width', percentage + '%');
        if (percentage == 100) {
            $('#progressbar .progress-bar').removeClass('progress-active');
            if (temporary) {
                $('#progressbar span').text($('#progressbar').data('old-text'));
            }
        }
    }
}

// Animates the search form from the middle of the page to the top right.
function animateSearchForm() {
    // A bit hackish: we rip the search form out of the DOM and use
    // position: absolute, so that we can later animate it across the page
    // into the top right #searchbox div.
    var sf = $('#searchform');
    var pos = sf.position();
    $('#searchbox .formplaceholder').css({ width: sf.width(), height: sf.height() });
    pos.position = 'absolute';
    $('#searchdiv .formplaceholder').css('height', sf.height());
    sf.detach();
    sf.appendTo('#content');
    sf.css(pos);

    sf.animate($('#searchbox').position(), 'fast', function() {
        $('#searchdiv').hide();
    });
}

function showResultsPage() {
    $('#results li').remove();
    $('#normalresults').show();
    $('#progressbar').show();
    $('#options').hide();
    $('#packageshint').hide();
    $('#pagination').text('');
    $('#perpackage-pagination').text('');
}

function sendQuery() {
    showResultsPage();
    $('#packages').text('');
    $('#errors div.alert-danger').remove();
    var query = {
        "Query": "q=" + encodeURIComponent(searchterm),
    };
    connection.send(JSON.stringify(query));
    document.title = searchterm + ' · Debian Code Search';
    var entries = localStorage.getItem("autocomplete");
    if (entries === null) {
        localStorage["autocomplete"] = JSON.stringify([searchterm]);
    } else {
        entries = JSON.parse(entries);
        if (entries.indexOf(searchterm) === -1) {
            entries.push(searchterm);
        }
        localStorage["autocomplete"] = JSON.stringify(entries);
    }
    animateSearchForm();

    progress(0, false, 'Checking which files to grep…');
}

connection.onopen = function() {
    clearTimeout(showConnectProgress);
    $('#searchform input').attr('disabled', false);

    // The URL dictates a search query, so start it.
    if (window.location.pathname.lastIndexOf('/results/', 0) === 0 ||
        window.location.pathname.lastIndexOf('/perpackage-results/', 0) === 0) {
        var parts = new RegExp("results/([^/]+)").exec(window.location.pathname);
        searchterm = decodeURIComponent(parts[1]);
        sendQuery();
    }

    $('#searchform').submit(function(ev) {
        searchterm = $('#searchform input[name=q]').val();
        sendQuery();
        history.pushState({ searchterm: searchterm, nr: 0, perpkg: false }, 'page ' + 0, '/results/' + encodeURIComponent(searchterm) + '/page_0');
        ev.preventDefault();
    });

    // This is triggered when the user navigates (e.g. via back button) between
    // pages that were created using history.pushState().
    $(window).bind("popstate", function(ev) {
        var state = ev.originalEvent.state;
        console.log('popstate', state);
        if (state == null) {
            // Restore the original page.
            $('#normalresults, #perpackage, #progressbar, #errors, #packages, #options').hide();
            $('#searchdiv').show();
            $('#searchdiv .formplaceholder').after($('#searchform'));
            $('#searchform').css('position', 'static');
            restoreAutocomplete();
        } else {
            if (!$('#normalresults').is(':visible') &&
                !$('#perpackage').is(':visible')) {
                showResultsPage();
                animateSearchForm();
                // The following are necessary because we don’t send the query
                // anew and don’t get any progress messages (the final progress
                // message triggers displaying certain elements).
                $('#packages, #errors, #options').show();
            }
            $('#enable-perpackage').prop('checked', state.perpkg);
            changeGrouping();
            if (state.perpkg) {
                loadPerPkgPage(state.nr);
            } else {
                loadPage(state.nr);
            }
        }
    });
};

connection.onerror = function(e) {
    // We could display an error, but since the page is supposed to fall back
    // gracefully, why would the user be concerned if the search takes a tiny
    // bit longer than usual?
    // error(false, true, 'websocket_broken', 'Could not open WebSocket connection to ' + e.target.URL);
};

connection.onclose = function(e) {
    // XXX: ideally, we’d only display the message if the reconnect takes longer than, say, a second?
    var msg = error(false, false, null, 'Lost connection to Debian Code Search. Reconnecting…');
    $('#searchform input').attr('disabled', true);

    var oldHandler = connection.onopen;
    connection.onopen = function() {
        $('#searchform input').attr('disabled', false);
        msg.remove();
        connection.onopen = oldHandler;
        oldHandler();
    };
};

var queryid;
var resultpages;
var currentpage;
var currentpage_pkg;
var packages = [];

function addSearchResult(results, result) {
    var context = [];
    // NB: All of the following context lines are already HTML-escaped by the server.
    context.push(result.Ctxp2);
    context.push(result.Ctxp1);
    context.push('<strong>' + result.Context + '</strong>');
    context.push(result.Ctxn1);
    context.push(result.Ctxn2);
    // Remove any empty context lines (e.g. when the match is close to the
    // beginning or end of the file).
    context = $.grep(context, function(elm, idx) { return $.trim(elm) != ""; });
    context = context.join("<br>").replace("\t", "    ");

    // Split the path into source package (bold) and rest.
    var delimiter = result.Path.indexOf("_");
    var sourcePackage = result.Path.substring(0, delimiter);
    var rest = result.Path.substring(delimiter);

    // Append the new search result, then sort the results.
    results.append('<li data-ranking="' + result.Ranking + '"><a href="/show?file=' + encodeURIComponent(result.Path) + '&line=' + result.Line + '"><code><strong>' + sourcePackage + '</strong>' + escapeForHTML(rest) + '</code></a><br><pre>' + context + '</pre><small>PathRank: ' + result.PathRank + ', Final: ' + result.Ranking + '</small></li>');
    $('ul#results').append($('ul#results>li').detach().sort(function(a, b) {
        return b.getAttribute('data-ranking') - a.getAttribute('data-ranking');
    }));

    // For performance reasons, we always keep the amount of displayed
    // results at 10. With (typically rather generic) queries where the top
    // results are changed very often, the page would get really slow
    // otherwise.
    var items = $('ul#results>li');
    if (items.size() > 10) {
        items.last().remove();
    }

    fixProgressbar();
}

function loadPage(nr) {
    // Start the progress bar after 20ms. If the page was in the cache, this
    // timer will be cancelled by the load callback below. If it wasn’t, 20ms
    // is short enough of a delay to not be noticed by the user.
    var progress_bar_start = setTimeout(function() {
        progress(0, true, 'Loading search result page ' + nr + '…');
    }, 20);

    var pathname = '/results/' + encodeURIComponent(searchterm) + '/page_' + nr;
    if (location.pathname != pathname) {
        history.pushState({ searchterm: searchterm, nr: nr, perpkg: false }, 'page ' + nr, pathname);
    }
    $.ajax('/results/' + queryid + '/page_' + nr + '.json')
        .done(function(data, textStatus, xhr) {
            clearTimeout(progress_bar_start);
            // TODO: experiment and see whether animating the results works
            // well. Fade them in one after the other, see:
            // http://www.google.com/design/spec/animation/meaningful-transitions.html#meaningful-transitions-hierarchical-timing
            currentpage = nr;
            updatePagination($('#pagination'), currentpage, resultpages, 'loadPage');
            $('ul#results>li').remove();
            var ul = $('ul#results');
            $.each(data, function(idx, element) {
                addSearchResult(ul, element);
            });
            progress(100, true, null);
        })
        .fail(function(xhr, textStatus, errorThrown) {
            error(true, true, null, 'Could not load search query results ("' + errorThrown + '").');
        });
}

// If preload is true, the current URL will not be updated, as the data is
// preloaded and inserted into (hidden) DOM elements.
function loadPerPkgPage(nr, preload) {
    var progress_bar_start;
    if (!preload) {
        // Start the progress bar after 20ms. If the page was in the cache,
        // this timer will be cancelled by the load callback below. If it
        // wasn’t, 20ms is short enough of a delay to not be noticed by the
        // user.
        progress_bar_start = setTimeout(function() {
            progress(0, true, 'Loading search result page ' + nr + '…');
        }, 20);
        var pathname = '/perpackage-results/' + encodeURIComponent(searchterm) + '/2/page_' + nr;
        if (location.pathname != pathname) {
            history.pushState({ searchterm: searchterm, nr: nr, perpkg: true }, 'page ' + nr, pathname);
        }
    }
    $.ajax('/perpackage-results/' + queryid + '/2/page_' + nr + '.json')
        .done(function(data, textStatus, xhr) {
            if (progress_bar_start !== undefined) {
                clearTimeout(progress_bar_start);
            }
            currentpage_pkg = nr;
            updatePagination($('#perpackage-pagination'), currentpage_pkg, Math.trunc(packages.length / packagesPerPage), 'loadPerPkgPage');
            var pp = $('#perpackage-results');
            pp.text('');
            $.each(data, function(idx, meta) {
                pp.append('<h2>' + meta.Package + '</h2>');
                var ul = $('<ul></ul>');
                pp.append(ul);
                $.each(meta.Results, function(idx, result) {
                    addSearchResult(ul, result);
                });
                if (!preload) {
                    progress(100, true, null);
                }
            });
        })
        .fail(function(xhr, textStatus, errorThrown) {
            error(true, true, null, 'Could not load search query results ("' + errorThrown + '").');
        });
}

function updatePagination(p, currentpage, resultpages, clickFunc) {
    p.text('');
    p.append('<strong>Pages:<strong> ');
    if (currentpage > 0) {
        p.append('<a href="javascript: ' + clickFunc + '(0);">1</a> ');
        p.append('<span>&lt;</span> ');
    }
    var start = Math.max(currentpage - 5, (currentpage > 0 ? 1 : 0));
    var end = Math.min((currentpage >= 5 ? currentpage + 5 : 10), resultpages);

    for (var i = start; i < end; i++) {
        //if (i < 3) {
        //    p.append('<link rel="prerender" href="/results/' + msg.QueryId + '/page_' + i + '.json">');
        //}
        p.append('<a style="' + (i == currentpage ? "font-weight: bold" : "") + '" href="javascript: ' + clickFunc + '(' + i + ');">' + (i + 1) + '</a> ');
    }

    if (currentpage < (resultpages-1)) {
        p.append('<span>&gt;</span> ');
    }

    if (end < resultpages) {
        p.append('… <a href="javascript: ' + clickFunc + '(' + (resultpages - 1) + ');">' + resultpages + '</a>');
    }
}

function escapeForHTML(input) {
    return $('<div/>').text(input).html();
}

connection.onmessage = function(e) {
    var msg = JSON.parse(e.data);
    switch (msg.Type) {
        case "progress":
        queryid = msg.QueryId;

        progress(((msg.FilesProcessed / msg.FilesTotal) * 90) + 10,
                 false,
                 msg.FilesProcessed + ' / ' + msg.FilesTotal + ' files grepped (' + msg.Results + ' results)');
        if (msg.FilesProcessed == msg.FilesTotal) {
            if (msg.Results === 0) {
                error(false, true, 'noresults', 'Your query “' + searchterm + '” had no results. Did you read the FAQ?');
            } else {
                $('#options').show();

                progress(100, false, msg.FilesTotal + ' files grepped (' + msg.Results + ' results)');

                // Request the results, but grouped by Debian source package.
                // Having these available means we can directly show them when the
                // user decides to switch to perpackage mode.
                loadPerPkgPage(0, true);

                $.ajax('/results/' + queryid + '/packages.json')
                    .done(function(data, textStatus, xhr) {
                        var p = $('#packages');
                        p.text('');
                        packages = data.Packages;
                        updatePagination($('#perpackage-pagination'), currentpage_pkg, Math.trunc(packages.length / packagesPerPage), 'loadPerPkgPage');
                        if (data.Packages.length === 1) {
                            p.append('All results from Debian source package <strong>' + data.Packages[0] + '</strong>');
                        } else if (data.Packages.length > 1) {
                            // We are limiting the amount of packages because
                            // some browsers (e.g. Chrome 40) will stop
                            // displaying text with “white-space: nowrap” once
                            // it becomes too long.
                            var packagesList = data.Packages.slice(0, 1000).join(', ');
                            p.append('<span><strong>Filter by package</strong>: ' + packagesList + '</span>');
                            if ($('#packages span:first-child').prop('scrollWidth') > p.width()) {
                                p.append('<span class="showhint"><a href="#" onclick="$(\'#packageshint\').show(); return false;">▾</a></span>');
                                $('#packageshint').text('');
                                $('#packageshint').append('To see all packages which contain results: <pre>curl -s http://' + window.location.host + '/results/' + queryid + '/packages.json | jq -r \'.Packages[]\'</pre>');
                            }

                            $('#enable-perpackage').attr('disabled', null);
                            $('label[for=enable-perpackage]').css('opacity', '1.0');

                            if (window.location.pathname.lastIndexOf('/perpackage-results/', 0) === 0) {
                                var parts = new RegExp("/perpackage-results/([^/]+)/2/page_([0-9]+)").exec(window.location.pathname);
                                $('#enable-perpackage').prop('checked', true);
                                changeGrouping();
                                loadPerPkgPage(parseInt(parts[2]));
                            }
                        }
                    })
                    .fail(function(xhr, textStatus, errorThrown) {
                        error(true, true, null, 'Loading search result package list failed ("' + errorThrown + '").');
                    });
            }
        }
        break;

        case "pagination":
        // Store the values in global variables for constructing URLs when the
        // user requests a different page.
        resultpages = msg.ResultPages;
        queryid = msg.QueryId;
        currentpage = 0;
        currentpage_pkg = 0;
        updatePagination($('#pagination'), currentpage, resultpages, 'loadPage');

        if (window.location.pathname.lastIndexOf('/results/', 0) === 0) {
            var parts = new RegExp("/results/([^/]+)/page_([0-9]+)").exec(window.location.pathname);
            loadPage(parseInt(parts[2]));
        }
        break;

        case "result":
        addSearchResult($('ul#results'), msg);
        break;

        case "error":
        if (msg.ErrorType == "backendunavailable") {
            error(false, true, msg.ErrorType, "The results may be incomplete, not all Debian Code Search servers are okay right now.");
        } else {
            error(msg.ErrorType);
        }
        break;

        default:
        throw new Error('Server sent unknown piece of data, type is "' + msg.Type);
    }
};

function setPositionAbsolute(selector) {
    var element = $(selector);
    var pos = element.position();
    pos.width = element.width();
    pos.height = element.height();
    pos.position = 'absolute';
    element.css(pos);
}

function setPositionStatic(selector) {
    $(selector).css({
        'position': 'static',
        'left': '',
        'top': '',
        'width': '',
        'height': ''});
}

// Switch between displaying all results and grouping search results by Debian
// source package.
function changeGrouping() {
    var ppelements = $('#perpackage');

    var currentPerPkg = ppelements.is(':visible');
    var shouldPerPkg = $('#enable-perpackage').prop('checked');
    if (currentPerPkg === shouldPerPkg) {
        return;
    }

    ppelements.data('hideAfterAnimation', !shouldPerPkg);

    if (currentPerPkg) {
            $('#perpkg').addClass('animation-reverse');
    } else {
            $('#perpkg').removeClass('animation-reverse');
            $('#perpkg').show();
    }

    if (shouldPerPkg) {
        ppelements.removeClass('animation-reverse');
        var pathname = '/perpackage-results/' + encodeURIComponent(searchterm) + '/2/page_' + currentpage_pkg;
        if (location.pathname != pathname) {
            history.pushState(
                { searchterm: searchterm, nr: currentpage_pkg, perpkg: true },
                'page ' + currentpage_pkg,
                pathname);
        }

        setPositionAbsolute('#footer');
        setPositionAbsolute('#normalresults');
        $('#perpackage').show();
    } else {
        ppelements.addClass('animation-reverse');
        var pathname = '/results/' + encodeURIComponent(searchterm) + '/page_' + currentpage;
        if (location.pathname != pathname) {
            history.pushState(
                { searchterm: searchterm, nr: currentpage, perpkg: false },
                'page ' + currentpage,
                pathname);
        }
        $('#normalresults').show();
        // For browsers that don’t support animations, we need to have a fallback.
        // The timer will be cancelled in the animationstart event handler.
        animationFallback = setTimeout(function() {
            $('#perpackage').hide();
            setPositionStatic('#footer, #normalresults');
        }, 50);
    }

    ppelements.removeClass('ppanimation');
    // Trigger a reflow, otherwise removing/adding the animation class does not
    // lead to restarting the animation.
    ppelements[0].offsetWidth = ppelements[0].offsetWidth;
    ppelements.addClass('ppanimation');
}

// Restore autocomplete from localstorage. This is necessary because the form
// never gets submitted (we intercept the submit event). All the alternatives
// are worse and have side-effects.
function restoreAutocomplete() {
    var entries = localStorage.getItem("autocomplete");
    if (entries !== null) {
        entries = JSON.parse(entries);
        var dataList = document.getElementById('autocomplete');
        $('datalist').empty();
        $.each(entries, function() {
            var option = document.createElement('option');
            option.value = this;
            dataList.appendChild(option);
        });
    }
}

// This function needs to be called every time a scrollbar can appear (any DOM
// changes!) or the size of the window is changed.
//
// This is because span.progressbar-front-text needs to be the same width as
// div#progressbar, but there is no way to specify that in pure CSS :|.
function fixProgressbar() {
    $('.progressbar-front-text').css('width', $('#progressbar').css('width'));
}

$(window).load(function() {
    // Try to restore autocomplete settings even before the connection is
    // established. If localStorage contains an entry, the user has used the
    // instant search at least once, so chances are she’ll use it again.
    restoreAutocomplete();

    // Pressing “/” anywhere on the page focuses the search field.
    $(document).keydown(function(e) {
        if (e.which == 191) {
            var q = $('#searchform input[name=q]');
            if (q.is(':focus')) {
                return;
            }
            q.focus();
            e.preventDefault();
        }
    });

    fixProgressbar();

    $(window).resize(fixProgressbar);

    function bindAnimationEvent(element, name, cb) {
        var prefixes = ["webkit", "MS", "moz", "o", ""];
        for (var i = 0; i < prefixes.length; i++) {
            if (i >= 3) {
                element.bind(prefixes[i] + name.toLowerCase(), cb);
            } else {
                element.bind(prefixes[i] + name, cb);
            }
        }
    }

    var ppresults = $('#perpackage');
    bindAnimationEvent(ppresults, 'AnimationStart', function(e) {
        clearTimeout(animationFallback);
    });
    bindAnimationEvent(ppresults, 'AnimationEnd',  function(e) {
            if (ppresults.data('hideAfterAnimation')) {
                    ppresults.hide();
                    setPositionStatic('#footer, #normalresults');
            } else {
                    $('#normalresults').hide();
            }
    });

    if (window.location.pathname.lastIndexOf('/results/', 0) === 0 ||
        window.location.pathname.lastIndexOf('/perpackage-results/', 0) === 0) {
        var parts = new RegExp("results/([^/]+)").exec(window.location.pathname);
        $('#searchform input[name=q]').val(decodeURIComponent(parts[1]));

        // If the websocket is not connected within 100ms, indicate progress.
        if (connection.readyState != WebSocket.OPEN) {
            $('#searchform input').attr('disabled', true);
            showConnectProgress = setTimeout(function() {
                $('#progressbar').show();
                fixProgressbar();
                progress(0, true, 'Connecting…');
            }, 100);
        }
    }
});

