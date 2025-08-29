/*
 * Copyright (c) 2025, WSO2 LLC.
 * Licensed under the Apache License, Version 2.0.
 */
package org.wso2.apk.integration.reporting;

import io.cucumber.plugin.ConcurrentEventListener;
import io.cucumber.plugin.event.EventPublisher;
import io.cucumber.plugin.event.Result;
import io.cucumber.plugin.event.Status;
import io.cucumber.plugin.event.TestCase;
import io.cucumber.plugin.event.TestCaseFinished;
import io.cucumber.plugin.event.TestRunFinished;

import java.net.URI;
import java.nio.file.Paths;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

/**
 * Prints a compact pass/fail summary table for Cucumber scenarios to the console.
 */
public class ConsoleSummaryTablePlugin implements ConcurrentEventListener {

    private static class Row {
        String feature;
        String scenario;
        Status status;
        Duration duration;
    }

    private final List<Row> rows = new CopyOnWriteArrayList<>();

    @Override
    public void setEventPublisher(EventPublisher publisher) {
        publisher.registerHandlerFor(TestCaseFinished.class, this::onFinished);
        publisher.registerHandlerFor(TestRunFinished.class, this::onRunFinished);
    }

    private void onFinished(TestCaseFinished event) {
        TestCase tc = event.getTestCase();
        Result result = event.getResult();
        Row r = new Row();
        r.feature = featureName(tc.getUri());
        r.scenario = tc.getName();
        r.status = result.getStatus();
        r.duration = result.getDuration() == null ? Duration.ZERO : result.getDuration();
        rows.add(r);
    }

    private String featureName(URI uri) {
        if (uri == null) return "";
        try {
            return Paths.get(uri).getFileName().toString();
        } catch (Exception e) {
            String s = uri.toString();
            int slash = s.lastIndexOf('/');
            return slash >= 0 ? s.substring(slash + 1) : s;
        }
    }

    private void onRunFinished(TestRunFinished event) {
        if (rows.isEmpty()) {
            System.out.println("No scenarios executed.");
            return;
        }

    // Sort: FAILED first, then others, PASSED last
    List<Row> sorted = new ArrayList<>(rows);
    sorted.sort(Comparator
        .comparingInt((Row r) -> statusWeight(r.status))
        .thenComparing(r -> r.feature)
        .thenComparing(r -> r.scenario));

        // Compute column widths
        int featureW = Math.max("Feature".length(), sorted.stream().map(r -> r.feature).mapToInt(String::length).max().orElse(7));
        int scenarioW = Math.max("Scenario".length(), sorted.stream().map(r -> r.scenario).mapToInt(String::length).max().orElse(8));
        int statusW = "Status".length();
        int durationW = "Duration".length();

    // header string built dynamically below when printing
        String sep = "+" + repeat('-', featureW + 2) + "+" + repeat('-', scenarioW + 2)
                + "+" + repeat('-', statusW + 2) + "+" + repeat('-', durationW + 2) + "+";

        System.out.println();
        System.out.println("Cucumber scenario summary:");
        System.out.println(sep);
        System.out.println(String.format("| %" + (-featureW) + "s | %" + (-scenarioW) + "s | %" + (-statusW) + "s | %" + (-durationW) + "s |",
                "Feature", "Scenario", "Status", "Duration"));
        System.out.println(sep);
        for (Row r : sorted) {
            String dur = prettyDuration(r.duration);
            System.out.println(String.format("| %" + (-featureW) + "s | %" + (-scenarioW) + "s | %" + (-statusW) + "s | %" + (-durationW) + "s |",
                    r.feature, r.scenario, r.status.name(), dur));
        }
        System.out.println(sep);

        long passed = sorted.stream().filter(r -> r.status == Status.PASSED).count();
        long failed = sorted.stream().filter(r -> r.status == Status.FAILED).count();
        long skipped = sorted.stream().filter(r -> r.status == Status.SKIPPED).count();
        System.out.printf("Totals: %d passed, %d failed, %d skipped, %d total\n",
                passed, failed, skipped, sorted.size());
        System.out.println();
    }

    private static String repeat(char c, int n) {
        StringBuilder sb = new StringBuilder(n);
        for (int i = 0; i < n; i++) sb.append(c);
        return sb.toString();
    }

    private static String prettyDuration(Duration d) {
        long ms = d.toMillis();
        if (ms < 1000) return ms + "ms";
        long s = ms / 1000; ms %= 1000;
        long m = s / 60; s %= 60;
        if (m > 0) return String.format("%dm %ds", m, s);
        return String.format("%ds %dms", s, ms);
    }

    private static int statusWeight(Status s) {
        if (s == Status.FAILED) return 0;
        // Group other non-passed statuses in the middle
        if (s == Status.AMBIGUOUS || s == Status.PENDING || s == Status.SKIPPED || s == Status.UNDEFINED) return 1;
        return 2; // PASSED last
    }
}
